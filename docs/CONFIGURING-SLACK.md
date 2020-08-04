# Configuring Slack

## Create Slack App
- Access https://api.slack.com/apps?new_app=1 in your workspace;
- Click on the button: __Create an App__;
- A model will be shown to you and you need to define your __App Name__ and select your __Slack Workspace__;
- Click on the button: __Create App__;
- In __Features__ click on the __OAuth & Permissions__;
- In __Scopes__/__Bot Token Scopes__ click on the __Add an OAuth Scope__;
- Then select the follows scopes:


### Bot Token Scopes
```
 - app_mentions:read
 - channels:join
 - channels:manage
 - channels:read
 - chat:write.public
 - chat:write
 - commands
 - pins:read
 - pins:write
 - usergroups:read
 - users:read
 - users:read.email
```

### User Token Scopes
```
- channels:read
```

## Install the app to workspace
- In __Features__ click on the __OAuth & Permissions__;
- And click on the __Install App to Workspace__;
- Select a channel for Slack to post as an app and click on the __Allow__;
- Copy the __Bot User OAuth Access Token__;
- Paste de code into the `HELLPER_OAUTH_TOKEN` variable;
- Restart your application.


## Events
- Now, in __Features__, click on the __Event Subscriptions__;
- And in __Enable Events__ turn on it;
- In __Request URL__, set your application's public URL to the field. It will look something like this: `https://yourhost.publicaddress.com/events`;
- Under it you should receive a request that looks like this:
```
Our Request:
POST
"body": {
	 "type": "url_verification",
	 "token": "XXXXXXXX",
	 "challenge": "YYYYYYY"
}
Your Response:
"code": 200
"error": "challenge_failed"
"body": {

}
```
_* you will receive the same info in your console application too._;

- Copy the `token` and past into the `HELLPER_VERIFICATION_TOKEN` variable;
- Restart your application, then the changes will be applied.


## Subscribe to bot events

- In the same page open the __Subscribe to bot events__, click on the __Add Bot User Event__ and add the `app_mention` option;
- After that, open the __Subscribe to events on behalf of users__, click on the __Add Workspace Event__ and add the `channel_created` option;
- Now, go in __Features__, click on the __Slash Commands__ and click on the __Create New Command__ to add the following commands:


| Command  | Request URL | Short Description |
| - | - | - |
|`/hellper_incident`|https://yourhost.publicaddress.com/open|_Starts Incident_|
|`/hellper_status`|https://yourhost.publicaddress.com/status|_Show all pinned messages_|
|`/hellper_close`|https://yourhost.publicaddress.com/close|_Closes Incident_|
|`/hellper_resolve`|https://yourhost.publicaddress.com/resolve|_Resolves Incident_|
|`/hellper_cancel`|https://yourhost.publicaddress.com/cancel|_Cancels Incident_|
|`/hellper_pause_notify`|https://yourhost.publicaddress.com/pause-notify|_Pauses incident notification_|
|`/hellper_update_dates`|https://yourhost.publicaddress.com/dates|_Updates the dates for an incident_|

- At last, click on the __Save Changes__;

- Now, in __Features__/__Interactivity & Shortcuts__ turn on the option __Interactivity__ and configure your address URL `http://yourhost.publicaddress.com/interactive`;
