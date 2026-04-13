import { UserLoginRequest, UserRegisterRequest } from "../api/requests/user";
import { createValidationScheme } from "../utils/create-validators-schema";
import * as Yup from "yup";

// Todo: i8n? 
export const userRegisterFormSchema = createValidationScheme<UserRegisterRequest>(
  Yup.object({
    // This does not allow accents. TODO: change to a more permissive regex
    name: Yup.string().matches(/^[a-zA-Z]+$/, "El nombre solamente puede contener letras").required(),
    surname: Yup.string().matches(/^[a-z A-Z]+$/, "El apellido solamente puede contener letras").required(),
    email: Yup.string().email("El correo electrónico debe ser válido").required(),
    password: Yup.string().min(6, "La contraseña debe tener al menos 6 caracteres").required(),
    confirmPassword: Yup.string().oneOf([Yup.ref('password'), ''], 'Las contraseñas deben coincidir').required(),
  })
);

export const usLoginFormSchema = createValidationScheme<UserLoginRequest>(
  Yup.object({
    email: Yup.string().email().required(),
    password: Yup.string().min(6).required(),
  })
)