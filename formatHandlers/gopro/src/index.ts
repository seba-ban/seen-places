import { startTransportHandling } from "./tansport";

startTransportHandling().catch((err) => {
  // TODO: use logger
  console.log("something went wrong");
  console.log(err);
  process.exit(1);
});
