import { createApp } from "@backstage/frontend-defaults";
import catalogPlugin from "@backstage/plugin-catalog/alpha";
import kubernetesPlugin from "@backstage/plugin-kubernetes/alpha"; // <- default import
import { navModule } from "./modules/nav";

export default createApp({
  features: [catalogPlugin, navModule, kubernetesPlugin],
});

// import { createApp } from "@backstage/frontend-defaults";
// import catalogPlugin from "@backstage/plugin-catalog/alpha";
// import { navModule } from "./modules/nav";

// export default createApp({
//   features: [catalogPlugin, navModule],
// });
