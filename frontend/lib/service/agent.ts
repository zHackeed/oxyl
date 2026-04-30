import { Caller } from '../api/api';
import { Agent } from '../api/models/agent';
import { AgentMetricsResponse } from '../api/models/metrics';
import { CreateAgentRequest } from '../api/requests/agent';
import { createAgentSchemaValidator } from '../validators/agent';

export const agentService = {
  get: async (companyId: string): Promise<Agent[] | null> => {
    try {
      const response = await Caller.get<Agent[]>(`/company/${companyId}/agents`);

      if (response.status !== 200) {
        console.error('Failed to fetch agents', response);
        return null;
      }

      return response.data;
    } catch (error) {
      console.error('Failed to fetch agents', error);
      return null;
    }
  },

  getOne: async (agentId: string): Promise<Agent | null> => {
    const response = await Caller.get(`/agent/${agentId}`);

    if (response.status !== 200) {
      console.error('Failed to fetch agent', response);
      return null;
    }

    console.log(response.data);

    return response.data as Agent;
  },

  create: async (agent: CreateAgentRequest): Promise<Agent | null> => {
    try {
      const valid = await createAgentSchemaValidator.validate(agent, {
        abortEarly: true,
      });

      if (!valid) {
        return Promise.reject('The schema is invalid!');
      }

      const response = await Caller.post<Agent>(`/agent/register`, agent);

      if (response.status !== 201) {
        console.error('Failed to create agent', response);
        return Promise.reject('Failed to create agent');
      }

      return response.data;
    } catch (error) {
      return Promise.reject(error);
    }
  },

  fetchAgentMetrics: async (
    agentId: string,
    interval: string
  ): Promise<AgentMetricsResponse | null> => {
    const response = await Caller.get<AgentMetricsResponse>(
      `/agent/${agentId}/metrics/${interval}`
    );

    if (response.status !== 200) {
      console.error('Failed to fetch agent metrics', response);
      return null;
    }

    return response.data;
  },

  fetchNotifications: async (agentId: string): Promise<any | null> => {
    const response = await Caller.get(`/agent/${agentId}/notifications`);
    
    if (response.status !== 200) {
      console.error('Failed to fetch agent notifications', response);
      return null;
    }

    return response.data;
  }
};
