export interface ActionRequest {
  request_id: string;
  action: string;
  target?: string;
  value?: string;
}

export interface ActionResponse {
  request_id: string;
  result: string;
  error?: string;
}

export interface ChatMessage {
  id: string;
  sender: 'user' | 'system' | 'ai';
  text: string;
}
