export interface Credentials {
  email: string;
  password: string;
}

export interface MeResponse {
  userId: string;
}

export interface TokenizeResponse {
  token: string;
  last4: string;
}

export interface DetokenizeResponse {
  pan: string;
}