type CreateCompanyRequest = {
  display_name: string;
  webhook_type: 'DISCORD' | 'SLACK';
  webhook_endpoint: string;
  webhook_channel?: string;
};

export { CreateCompanyRequest };
