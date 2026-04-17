import { createValidationScheme } from '../utils/create-validators-schema';
import { CreateAgentRequest } from '../api/requests/agent';
import * as Yup from 'yup';

export const createAgentSchemaValidator = createValidationScheme<CreateAgentRequest>(
  Yup.object({
    holder: Yup.string()
      .max(255, 'El nombre del agente no puede exceder 255 caracteres')
      .required('El nombre del agente es requerido'),
    display_name: Yup.string()
      .max(255, 'El nombre del agente no puede exceder 255 caracteres')
      .required('El nombre del agente es requerido'),
    registered_ip: Yup.string()
      .matches(
        /^(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]\d|\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]\d|\d)){3}$/,
        'La IP del agente debe ser válida'
      ) // https://github.com/sindresorhus/ip-regex/blob/main/index.js#L8
      .required('La IP del agente es requerida'),
  }).noUnknown(true)
);
