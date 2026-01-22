import Toybox.Test;
import Toybox.Lang;

module WhereAreYou {
    (:test)
    class RequestManagerTest {
        
        (:test)
        function testValidateInput(logger as Logger) as Boolean {
            var manager = new RequestManager("http://localhost");
            
            // Test Valid
            if (!manager.validateInput("user1")) {
                logger.debug("Failed to validate valid input 'user1'");
                return false;
            }

            // Test Invalid (Empty)
            if (manager.validateInput("")) {
                logger.debug("Failed to reject empty input");
                return false;
            }

            return true;
        }
    }
}
