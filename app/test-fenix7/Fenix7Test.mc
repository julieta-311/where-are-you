import Toybox.Test;
import Toybox.Lang;

(:test)
function testFenix7Module(logger as Logger) as Boolean {
    var metrics = new DeviceModule.Fenix7Metrics();
    var info = metrics.getDeviceSpecificInfo();
    
    if (info.equals("Fenix 7 Series Metrics Active")) {
        return true;
    }
    
    logger.debug("Unexpected info: " + info);
    return false;
}
