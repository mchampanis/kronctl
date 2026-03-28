// Paste this into Chrome DevTools console BEFORE connecting the keyboard
// in Keychron Launcher (launcher.keychron.com).
//
// Intercepts all WebHID sendReport / sendFeatureReport calls and logs
// non-keepalive packets to the console.

const origSendReport = HIDDevice.prototype.sendReport;
const origSendFeature = HIDDevice.prototype.sendFeatureReport;

HIDDevice.prototype.sendReport = function (reportId, data) {
  const arr = new Uint8Array(data);
  // 0xa3 is a keepalive/poll packet -- suppress to reduce noise
  if (arr[0] !== 0xa3) {
    console.log("sendReport", reportId, arr);
  }
  return origSendReport.call(this, reportId, data);
};

HIDDevice.prototype.sendFeatureReport = function (reportId, data) {
  const arr = new Uint8Array(data);
  if (arr[0] !== 0xa3) {
    console.log("sendFeatureReport", reportId, arr);
  }
  return origSendFeature.call(this, reportId, data);
};
