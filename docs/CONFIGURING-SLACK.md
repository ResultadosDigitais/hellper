# Configuring Slack

## Create Slack App

- Access <https://api.slack.com/apps>;
- Click on the button: __Create New App__;
- A model will be shown to you and you need to define your __App Name__ and select your __Slack Workspace__;
- Click on the button: __Create App__;

## Signing Secret

- In __Settings__/__App Credentials__, copy your __Signing Secret__;
- Paste it into the `HELLPER_SLACK_SIGNING_SECRET` variable;

## OAuth & Permissions - Bot Token Scopes

- In __Features__ click on the __OAuth & Permissions__;
- In __Scopes__/__Bot Token Scopes__ click on the __Add an OAuth Scope__;
- Then select the follows scopes:

```text
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

## Slash Commands

- Now, go in __Features__, click on the __Slash Commands__ and click on the __Create New Command__ to add the following commands:

| Command  | Request URL | Short Description |
| - | - | - |
|`/hellper_incident`|<https://yourhost.publicaddress.com/open>|_Starts Incident_|
|`/hellper_status`|<https://yourhost.publicaddress.com/status>|_Show all pinned messages_|
|`/hellper_close`|<https://yourhost.publicaddress.com/close>|_Closes Incident_|
|`/hellper_resolve`|<https://yourhost.publicaddress.com/resolve>|_Resolves Incident_|
|`/hellper_cancel`|<https://yourhost.publicaddress.com/cancel>|_Cancels Incident_|
|`/hellper_pause_notify`|<https://yourhost.publicaddress.com/pause-notify>|_Pauses incident notification_|
|`/hellper_update_dates`|<https://yourhost.publicaddress.com/dates>|_Updates the dates for an incident_|

## Interactivity & Shortcuts

- Now, in __Features__/__Interactivity & Shortcuts__ turn on the option __Interactivity__ and configure your address URL `http://yourhost.publicaddress.com/interactive`;

## Event Subscriptions

_Before that you need to start the Hellper application http server with the variable: `HELLPER_SLACK_SIGNING_SECRET`._

- Now, in __Features__, click on __Event Subscriptions__;
- And in __Enable Events__ turn on it;
- In __Request URL__, set your application's public URL to the field. It will look something like this: `https://yourhost.publicaddress.com/events`;
- In the same page open the __Subscribe to bot events__, click on the __Add Bot User Event__ and add the `app_mention` option;
- Click on __Save Changes__;

## OAuth Access Token

- In __Features__ click on __OAuth & Permissions__;
- Click on __Install App to Workspace__ and then click to __Allow__;
- Copy the __Bot User OAuth Access Token__;
- Paste the code into the `HELLPER_OAUTH_TOKEN` variable.
