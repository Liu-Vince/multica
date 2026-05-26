"use client";

import { useConfigStore } from "@multica/core/config";

export function useInstallCommands() {
  const installScriptUrl = useConfigStore((s) => s.installScriptUrl);
  const serverUrl = useConfigStore((s) => s.serverUrl);
  const appUrl = useConfigStore((s) => s.appUrl);

  const INSTALL_CMD = `curl -fsSL ${installScriptUrl} | bash`;

  const SETUP_CMD = (() => {
    if (serverUrl && appUrl) {
      return `multica setup self-host --server-url ${serverUrl} --app-url ${appUrl}`;
    }
    return "multica setup";
  })();

  const TOKEN_CMD = `multica config set server_url ${serverUrl || "https://api.multica.ai"}
multica config set app_url ${appUrl || "https://multica.ai"}
multica login --token <YOUR_TOKEN>
multica daemon start`;

  return { INSTALL_CMD, SETUP_CMD, TOKEN_CMD } as const;
}
