export type InputType = "Text" | "Image" | "ImageHard";

export interface SubmitRequest {
  type: "submit_request";
  request_id: string;
  input_type: InputType;
  content?: string | null;
  image_base64?: string | null;
  metadata?: Record<string, any> | null;
}

export interface CancelRequest {
  type: "cancel_request";
  request_id: string;
}

export type ClientMessage = SubmitRequest | CancelRequest;

export interface AckMessage {
  type: "ack";
  request_id: string;
  status: string;
}

export interface ProgressMessage {
  type: "progress";
  request_id: string;
  stage: string;
  message: string;
}

export interface TokenMessage {
  type: "token";
  request_id: string;
  token: string;
  is_final_token: boolean;
}

export interface FinalResponseMessage {
  type: "final_response";
  request_id: string;
  response_text: string;
  model_used?: string;
}

export interface ErrorMessage {
  type: "error";
  request_id?: string;
  error_code?: string;
  message: string;
}

export type ServerMessage =
  | AckMessage
  | ProgressMessage
  | TokenMessage
  | FinalResponseMessage
  | ErrorMessage;

export type ChatRole = "user" | "bot" | "system";

export interface ChatMessage {
  id: string;
  role: ChatRole;
  text?: string;
  imageSrc?: string;
  partial?: boolean;
}

export {};
