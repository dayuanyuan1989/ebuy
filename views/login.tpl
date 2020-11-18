<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
    <style>

        body{
            //background-color: #00FFFF;
        }

    </style>
</head>
<body>
    <form action="" method="post">
        用户名:<input type="text" name="username"/><br>
        密码:<input type="password"  name="password"/><br>
        <input type="submit" value="登录">


    </form><br>
    <div>


        <a href="{{.website}}">{{.website}}}</a>
    </div>
    <div>
        {{if .isDisplay}}
            <em>{{.content1}}</em>
        {{else}}
            <em>{{.content2}}</em>
        {{end}}
        <h2>Users表</h2>
        <table border="1">

        <tr><th>Id</th><th>username</th><th>Password</th></tr>



        {{range .users}}
            <tr> <td>{{.Id}}</td><td>{{.Username}}</td><td>{{.Password}}</td></tr>

        {{end}}
        </table>




    </div>
</body>
</html>