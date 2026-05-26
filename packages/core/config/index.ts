import { createStore } from "zustand/vanilla";
import { useStore } from "zustand";

export interface CLIConfig {
  installScriptUrl: string;
  cliDownloadBaseUrl: string;
  serverUrl: string;
  appUrl: string;
}

interface ConfigState {
  cdnDomain: string;
  allowSignup: boolean;
  googleClientId: string;
  /** CLI / install URLs from the backend /api/config. Defaults match the
   *  public GitHub / Multica Cloud values so existing deployments keep
   *  working without setting new env vars. */
  installScriptUrl: string;
  cliDownloadBaseUrl: string;
  serverUrl: string;
  appUrl: string;
  setCdnDomain: (domain: string) => void;
  setAuthConfig: (config: {
    allowSignup: boolean;
    googleClientId?: string;
  }) => void;
  setCLIConfig: (config: Partial<CLIConfig>) => void;
}

const DEFAULT_INSTALL_SCRIPT_URL =
  "https://raw.githubusercontent.com/multica-ai/multica/main/scripts/install.sh";
const DEFAULT_CLI_DOWNLOAD_BASE_URL =
  "https://github.com/multica-ai/multica/releases/latest/download";

export const configStore = createStore<ConfigState>((set) => ({
  cdnDomain: "",
  allowSignup: true,
  googleClientId: "",
  installScriptUrl: DEFAULT_INSTALL_SCRIPT_URL,
  cliDownloadBaseUrl: DEFAULT_CLI_DOWNLOAD_BASE_URL,
  serverUrl: "",
  appUrl: "",
  setCdnDomain: (domain) => set({ cdnDomain: domain }),
  setAuthConfig: ({ allowSignup, googleClientId = "" }) =>
    set({ allowSignup, googleClientId }),
  setCLIConfig: (config) =>
    set({
      installScriptUrl:
        config.installScriptUrl ?? DEFAULT_INSTALL_SCRIPT_URL,
      cliDownloadBaseUrl:
        config.cliDownloadBaseUrl ?? DEFAULT_CLI_DOWNLOAD_BASE_URL,
      serverUrl: config.serverUrl ?? "",
      appUrl: config.appUrl ?? "",
    }),
}));

export function useConfigStore(): ConfigState;
export function useConfigStore<T>(selector: (state: ConfigState) => T): T;
export function useConfigStore<T>(selector?: (state: ConfigState) => T) {
  return useStore(configStore, selector as (state: ConfigState) => T);
}