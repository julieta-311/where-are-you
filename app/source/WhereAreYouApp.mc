import Toybox.Application;
import Toybox.Lang;
import Toybox.WatchUi;

class WhereAreYouApp extends Application.AppBase {
    var _requestManager as WhereAreYou.RequestManager?;

    function initialize() {
        AppBase.initialize();
    }

    function onStart(state as Dictionary?) as Void {
        // In simulator, localhost is 127.0.0.1.
        _requestManager = new WhereAreYou.RequestManager("http://127.0.0.1:8080/api/v1");
    }

    function onStop(state as Dictionary?) as Void {
    }

    function getInitialView() as [WatchUi.Views] or [WatchUi.Views, WatchUi.InputDelegates] {
        if (_requestManager != null) {
            return [ new WhereAreYouView(), new WhereAreYou.WhereAreYouDelegate(_requestManager) ];
        }
        return [ new WhereAreYouView() ];
    }
}

