# Documentación Técnica

## Arquitectura General

El sistema se compone de varios componentes:
- Un paquete definido como `shared`, en el cual principalmente se encuentran todas los componentes comunes del backend (y sus módulos). Siendo en su composición
  - Estructuras tipadas de los modelos
  - Conectores de las diferentes bases de datos (tanto psql y redis). Con las correspondientes adaptaciones para su uso
  - Utilidades compartidas (obtenciones de informacion del contexto, manejo de errores, etc.)
  - Estructuras de manejo de datos (crud y sus necesidades)
  - Servicios correspondientes para las diferentes funciones generales
  > **Nota:** En este principalmente se realizan validaciones y autenticacion de permisos. Etc...

Esto permite que se pueda reutilizar de forma segura y consistente en todos los servicios. Evitando duplicacion.

- Módulo API - primer servicio encargado de todas las operaciones CRUD y modificaciones. Siendo el entrypoint del sistema.
  - Esta basado en Fiber con una arquitectura modular para poder expandirse rápidamente (en el cual es de apoyo el paquete shared)
  - Es el punto de autenticacion para tanto los usuarios como los agentes (será explicado mas adelante)

- Módulo Protocol - es el modulo general compartido donde estan las definiciones de los protocolos utilizados por los agentes y el ingress. Por lo general no se modifica. 

- Módulo Ingress - es el modulo encargado de la recepcion de los mensajes de los agentes. Se encarga de la recepcion de los mensajes y su procesamiento. 
  - Su comunicación es basada a través de Protocol Buffers, el cual su esquema está definido en `protocol`.
  - La autenticación del agente esta basada en 3 partes:
    - La autenticación inicial esta basada en `JWT`, la cual periodicamente se va renovando a través del modulo de API. Debido a las limitaciones de como es GRPC es más complejo el como mantener la sesión. 

    > **Nota:** Es posible que se pueda aceptar directamente y que se genere la autenticacion por parte inicial de la conexion, pero es complejidad añadida que habria que mantener en cuestion de como es el flujo de conexión.

    - Una token de autenticación que se genera al iniciar el agente por primera vez y este se encuentra en la fase de `ENROLL`. Sin ello, el agente no podrá realizar ningún tipo de publicación a los canales de monitorización.

  Hay que tomar en cuenta que el sistema una vez haya generado la token de autenticación, esta se almacenara en el agente para su uso posterior. Y si se pierde, el sistema no podrá autenticar al agente y debera ser reenrollado.

  En todo caso, tanto las autenticaciones como los tokens de autenticación y las claves de enroll estan generadas en base a una key de tipo ed25519. La cual se genera en la primera vez que se inicia el sistema y se almacena en su volumen de datos. (( En el cual se puede ver todo el flujo de inicio y explicacion en el README.md ))

- Modulo Agent - es como podemos menciona el córazon del sistema, este es el programa que se ejecuta en los sistemas los cuales se desea monitorear. Este se encarga de la ejecucion de las tareas y la comunicacion con el ingress. 
  - Como ya fue indicado, este solamente soporta sistemas basados en Linux, su adaptación para otros sistema de forma mas abstracta se puede llegar a hacer. Pero en su actualidad esta usando `procfs` el cual depende de toda la estructura del sistema de `/proc, /sys, /dev...` que usa Linux    

## Flujos de autenticación

### Autenticación de usuarios

Los usuarios principalmente llamarán a la API para obtener sus tokens de autenticación. Siendo correspondiente dos pares de tokens:
  - Una token de acceso, con una duracion corta. (El cual deberia ser 15 mins, aunque esta cambiado en su actualidad para las pruebas de que los diferentes tipos de endpoints funcionan)
  - Una token de refresco, con una duracion larga. (30 Dias). 

En todo caso, se intentará mantener un flujo de autenticacion que no se mantengan ningun de estado dentro del sistema, en su excepcion de las tokens de refresco. Estas debido a que tienen un identificador único que se usa para su invalidacion en 2 circunstancias:
  - Cierre de sesion del usuario
  - Expiracion de la token de acceso

Por como es JWT, es imposible de bloquear token de acceso. Lo cual es aceptable en esta circustancia.

## Autenticación de agents

Los agentes tienen principalmente 2 fases de autenticacion. Las cuales hemos mencionado arriba, pero definiremos mas a fondo:

#### Fase inicial

En esta situacion, no sabemos en que estado esta el agente. Ya que esta recien iniciado y no tiene ningun tipo de token de autenticacion guardada.

El cual, este iniciara su proceso de autenticacion con la API, enviando su token de identificador del sistema. La cual sera validada por la API y si es valida, se le devolvera un token de autenticacion que sera guardada en su ruta `/etc/oxyl`.

En esta fase, comprobaremos si:
  - El agente existe con ese identificador
  - Si esta marcado como activo (aun no esta completamente diseñado)
  - Si la IP en la cual estamos recibiendo la solicitud es la misma que tiene el agente registrado en nuestro lado

Si todas estas condiciones se cumplen, se le devolvera un token de autenticacion `JWT`, la cual guardara de forma ephemera en memoria. 

Este servicio que hemos mencionado, internamente se encargará del refresco para evitar fallos de autenticacion.

#### Fase de enroll

La fase de enroll es la segunda fase de autenticacion, la cual se ejecuta cuando el agente ya tiene un token de autenticacion guardada en su sistema. Esto le permite ya poder crear una conexion con el ingress.

Dependiendo de como esté el estado el sistema podrá pasar 2 cosas:
  - Se iniciará el proceso de enroll.
  - Se cargará la token de enroll registrada en el sistema.

En todo caso, lo que importa es en nuestra fase incial `Enrolling`. En la cual recibiremos todas la información del sistema y enriqueceremos nuestra base de datos para así poder mostrás mas adelante.

En la parte del ingress. La token que se generará está compuesta sobre todos los datos que recibe del servidor y la firmará con la token de enroll que se generó en el proceso de enroll.

> Cabe recalcar que ahora mismo no se está validando. Siendo que ahora mismo en ingress solamente se está cacheando la token (ya que esta si tenemos que sacarla de la información del agente) y no estamos comprobando signaturas.

Una vez haya pasado esta fase, el agente empezará de forma periodica a enviar información al ingress para que pueda ir enriqueciendo nuestra base de datos de timeseries.


### ¿Que se consume del agente?

El agente principalmente enviará los correspondientes datos:
- Uptime del sistema
- Uso de CPU medio (no por nucleo, aunque se puede expandir a por nucleo de forma sencilla)
- Consumo de RAM
- Consumo de disco
- Temperatura del sistema
- Uso de red (entrada y salida)

Todo esto, a través de diferentes librerias como:
- `procfs`
- `smartctl`

Nos permite facilmente obtener toda la información necesaria para nuestro sistema. (( Aunque tuvo un poco de problemas con `smartctl` en algunos sistemas. Pero fue cuestion de adaptarlo. ))

## Frontend

El frontend esta basado para dispositivos IOS (principalmente para iPhone), en la cual a través de React Native y Tamagui se presentará una interfaz de estilo `minimalista`. 

Este depende de 2 modulos:
- `API`: La cual es la base de toda la comunicación con el backend.
- `Websocket`: La cual cual ofrece metricas en tiempo vivo (Actualmente no definido a su completo, se comentará en STRUGGLES.md).

