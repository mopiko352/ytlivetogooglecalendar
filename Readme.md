# 概要
youtubeの予定をグーグルカレンダーにいれるやつ

# 前提条件
- gcloudのSetup
- GCPのSetup
- Dockerのインストール
- googleカレンダー作っておく

# デプロイ方法
- ${GCPのprojectid}-live-to-calendar-srcs　という名前のバケットを作っておく

deployments/terraform/dev/variables.tfのfunctionにいれたいチャンネルの情報入れる
function.configuration配列ふやしたら複数人分実行できる

- channelid　youtubeのチャンネルID
- name CloudFunctionsの名前

terraform applyする
```
$ docker run -it -v C:\Users\user\go\ytlivetogooglecalneder:/go/ytlivegooglecalender -v C:\Users\user\AppData\Roaming\gcloud\application_default_credentials.json:/root/.config/gcloud/application_default_credentials.json  -w /go/ytlivegooglecalender/deployments/terraform --entrypoint /bin/sh hashicorp/terraform:1.5

# cd dev
# terraform apply --var-file=variables.tfvars
```

variables.tfにいれたGCPのSecret Manager名と同じ名前のSecetをつくる
- ytlive-to-calendar-sa-key
terraform applyしたときにできたサービスアカウントのJSONKEYをそのままいれる
```
{
  "type": "service_account",
  "project_id": "",
  ...
}
```
というやつ
- carol-calendar-id
GOOGLEカレンダーのカレンダーIDをいれる
```
なんか文字@group.calendar.google.com
```

