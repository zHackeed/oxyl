import { createValidationScheme } from '../utils/create-validators-schema';
import { CreateCompanyRequest } from '../api/requests/company';
import * as Yup from 'yup';

export const createCompanySchemaValidator = createValidationScheme<CreateCompanyRequest>(
  Yup.object({
    display_name: Yup.string()
      .max(255, 'El nombre de la compañía no puede exceder 255 caracteres')
      .required('El nombre de la compañía es requerido'),
    webhook_type: Yup.string()
      .oneOf(['DISCORD', 'SLACK'], 'El tipo de webhook debe ser DISCORD o SLACK')
      .required('El tipo de webhook es requerido'),
    webhook_endpoint: Yup.string()
      .url('La URL del webhook debe ser una URL válida')
      .required('La URL del webhook es requerida'),
    webhook_channel: Yup.string()
      .max(255, 'El canal no puede exceder 255 caracteres')
      .optional(),
  })
);
