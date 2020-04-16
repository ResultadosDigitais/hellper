# Configuring Google
This instructions guide will give your application permission to make a copy of the post-mortem doc template in your Google Drive when a new incident is started.


## Contents
1. [Official documentation](#Official-documentation)
2. [Get a Client ID and Client Secret](#Get-a-Client-ID-and-Client-Secret)
   * [OAuth consent screen](#OAuth-consent-screen)
   * [Credentials](#Credentials)
3. [Authorization code](#Authorization-code)
4. [Signing into the application](#Signing-into-the-application)
5. [Enabling Google Drive API](#Enabling-Google-Drive-API)
6. [Template Post-mortem](#Template-Post-mortem)
7. [Setting environment variables](#Setting-environment-variables)


## Official documentation
The instructions here are a summary of the [official documentation](https://cloud.google.com/iap/docs/authentication-howto#authenticating_from_a_desktop_app).


## Get a Client ID and Client Secret
1. Open the [Google API Console Credentials](https://console.developers.google.com/apis/credentials) page.
2. From the project drop-down menu, select an existing project or create a new one.


### OAuth consent screen
1. On the **OAuth consent screen** page, select an **User Type**, then click **Create**.
2. Give a name for it and save it.


### Credentials
1. On the **Credentials** page, select **Create credentials**, then select **OAuth client ID**.
2. Under **Application type**, choose **Other** and give it a name, _ie. Hellper_. Then, click **Create**.
3. On this next page, take note of the **Client ID** and **Client Secret**. You'll need these going forward. Then, click **Ok**.
4. At last, click in the Download icon button to configure your credentials.
5. Copy the content file and paste it in your environment variable called: `HELLPER_GOOGLE_DRIVE_CREDENTIALS`.


## Authorization code
1. Copy the **Client ID** you got in the last step (_[Credentials](#credentials)_).
2. In the following URL, change `YOUR_CLIENT_ID_HERE` with the content from **Client ID**:

```
https://accounts.google.com/o/oauth2/v2/auth?client_id=YOUR_CLIENT_ID_HERE&response_type=code&scope=https://www.googleapis.com/auth/drive&access_type=offline&redirect_uri=urn:ietf:wg:oauth:2.0:oob
```

3. Access the new URL in your web browser.
4. Allow the permissions to your application to be able to access yours files.
5. On this next page, take note of the **Code**. You'll need this going forward.


## Signing into the application
1. Now you need to copy the **Client ID**, **Client Secret** and **Code** of the last steps (_[Credentials](#credentials) and [Authorization Code](#authorization-code)_), and replace them respectively in the follow command:

```shell
curl --data client_id="YOUR_CLIENT_ID_HERE" \
  --data client_secret="YOUR_CLIENT_SECRET_HERE" \
  --data code="YOUR_AUTH_CODE_HERE" \
  --data redirect_uri=urn:ietf:wg:oauth:2.0:oob \
  --data grant_type=authorization_code \
  https://oauth2.googleapis.com/token
```

2. Run it in your terminal and copy the response.
3. Past it in your environment variable called: `HELLPER_GOOGLE_DRIVE_TOKEN`.


### Example
**Run it in terminal**
```shell
curl --data client_id="YOUR_CLIENT_ID" \
  --data client_secret="YOUR_CLIENT_SECRET" \
  --data code="YOUR_AUTH_CODE" \
  --data redirect_uri=urn:ietf:wg:oauth:2.0:oob \
  --data grant_type=authorization_code \
  https://oauth2.googleapis.com/token
```

**Response**
```http
{
  "access_token": "xxxxxxxxxxxxxxxxxx",
  "expires_in": 3599,
  "refresh_token": "xxxxxxxxxxxxxxxxxx",
  "scope": "https://www.googleapis.com/auth/drive",
  "token_type": "Bearer"
}
```


## Enabling Google Drive API
Access [API Library](https://console.developers.google.com/apis/library/drive.googleapis.com), then click **Enable**.


## Template Post-mortem
1. Create new a file in your Google Doc and copy the ID from the file, like this:

`https://docs.google.com/document/d/YOUR_FILE_ID_IS_HERE/edit`

2. Paste the ID in your environment variable called: `HELLPER_GOOGLE_DRIVE_FILE_ID`.


## Setting environment variables
Now you need to change these three variables:

| Variable | Explanation |
| --- | --- |
|**HELLPER_GOOGLE_DRIVE_CREDENTIALS** |[Google Drive Credentials](/docs/CONFIGURING-GOOGLE.md#Get-a-Client-ID-and-Client-Secret)|
|**HELLPER_GOOGLE_DRIVE_FILE_ID**|[Google Drive FileId](/docs/CONFIGURING-GOOGLE.md#Template-Post-mortem) to your post-mortem template|
|**HELLPER_GOOGLE_DRIVE_TOKEN**|[Google Drive Token](/docs/CONFIGURING-GOOGLE.md#Signing-in-to-the-application)|
