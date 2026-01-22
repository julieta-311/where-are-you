import Toybox.Lang;
import Toybox.WatchUi;
import Toybox.System;

module WhereAreYou {
    class WhereAreYouDelegate extends WatchUi.BehaviorDelegate {
        private var _requestManager as RequestManager;

        function initialize(manager as RequestManager) {
            BehaviorDelegate.initialize();
            _requestManager = manager;
        }

        function onSelect() as Boolean {
            System.println("Select button pressed. Sending request...");
            // In a real app, you'd pick a friend from a list.
            // Here we hardcode for the prototype.
            _requestManager.createRequest("watch-user", "target-friend");
            return true;
        }

        function onMenu() as Boolean {
            return true;
        }
    }
}
