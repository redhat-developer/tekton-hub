declare global {
  interface Window {
    config: any;
  }
}

window.config = window.config || {
  API_URL: 'no API_URL  set',
  GH_CLIENT_ID: 'no GH_CLIENT_ID set',
};

export const ENTER_KEY = 13;
export const API_URL =
  'https://api-tekton-hub.apps.cluster-blr-f235.blr-f235.example.opentlc.com';
export const GH_CLIENT_ID = window.config.GH_CLIENT_ID;

console.debug(`Using config: ${API_URL} | config:`, window.config);
