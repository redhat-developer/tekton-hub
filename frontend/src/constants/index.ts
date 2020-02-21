export const ENTER_KEY = 13;
export const API_URL = window.config.API_URL;
export const GH_CLIENT_ID = window.config.GH_CLIENT_ID;

declare global {
  interface Window {
    config: any;
  }
}

console.debug(`Using config: ${API_URL} | config:`, window.config);
