import Toybox.Communications;
import Toybox.Lang;
import Toybox.System;

module WhereAreYou {
    // RequestManager handles all HTTP communications with the backend.
    class RequestManager {
        private var _baseUrl as String;

        // initialize sets the base URL for the API.
        // For simulator, use http://localhost:8080/api/v1
        // For real device, use your computer's IP or public server.
        function initialize(baseUrl as String) {
            _baseUrl = baseUrl;
        }

        // validateInput checks if the IDs are valid (not empty).
        public function validateInput(id as String) as Boolean {
            if (id == null || id.length() == 0) {
                return false;
            }
            return true;
        }

        // createRequest sends a location request to a target user.
        public function createRequest(requesterId as String, targetId as String) as Void {
            if (!validateInput(requesterId) || !validateInput(targetId)) {
                System.println("Invalid Input: IDs cannot be empty.");
                return;
            }

            var url = _baseUrl + "/requests";
            var params = {
                "requester_id" => requesterId,
                "target_id" => targetId
            };
            
            var options = {
                :method => Communications.HTTP_REQUEST_METHOD_POST,
                :headers => {
                    "Content-Type" => Communications.REQUEST_CONTENT_TYPE_JSON
                },
                :responseType => Communications.HTTP_RESPONSE_CONTENT_TYPE_JSON
            };

            Communications.makeWebRequest(url, params, options, method(:onRequestCompleted));
            System.println("Sending request to: " + url);
        }

        // onRequestCompleted handles the HTTP response.
        function onRequestCompleted(responseCode as Number, data as Dictionary?) as Void {
            if (responseCode == 201) {
                System.println("Request Created Successfully: " + data);
            } else {
                System.println("Request Failed. Code: " + responseCode);
                if (data != null) {
                    System.println("Error Data: " + data);
                }
            }
        }
    }
}
