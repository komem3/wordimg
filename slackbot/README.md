# wordimg slackbot

It is slack command that generates emoji.

![image](https://user-images.githubusercontent.com/29977786/107849016-a8db4400-6e3b-11eb-8a36-be7c6a0162a4.png)

## Install Commands

1. Deploy to cloud funciton.

```shell
./deploy.sh $your_projectID
```

2. Set the `URL` environment variable to the URL of wordimg. (https://wordimg-otho5yxlgq-an.a.run.app/wordimg)

3. Creat slack app.

4. Add slash commands. URL is the url of deployed cloud funciton.

   - /wordimg
   - /wordimg1
   - /wordimg2
   - /wordimg3

5. Add a `chat:write` to oauth scope.

6. Install slack app.

7. Create secret.

   - `SLACK_SECRET` Signing Secret of slack.
   - `SLACK_APIKEY` Bot User OAuth Access Token of slack.

8. Grant secret accessor privileges to the cloud function service account.
