import { createValidationScheme } from "../utils/create-validators-schema";
import { CreateCompanyRequest } from "../api/requests/company";
import * as Yup from "yup";

export const createCompanySchemaValidator = createValidationScheme<CreateCompanyRequest>(
  Yup.object({
    display_name: Yup
      .string()
      .max(255, "El nombre de la compañía no puede exceder 255 caracteres")
      .required("El nombre de la compañía es requerido"),
  })
);