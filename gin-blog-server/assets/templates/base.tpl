{{define "base"}}
    <!DOCTYPE html>
    <html>
    <head>
        <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        <meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>
        {{template "styles" .}}
        <title>{{ .Subject}}</title>
    </head>
    <body class="">
    <table role="presentation" border="0" cellpadding="0" cellspacing="0" class="body">
        <tr>
            <td>&nbsp;</td>
            <td class="container">
                <div class="header">
                    <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                            <td class="align-center" width="100%">
                                我们共享过的月色，如今还在天上
                            </td>
                        </tr>
                    </table>
                </div>
                <div class="content">
                    <!-- START CENTERED WHITE CONTAINER -->
                    <span class="preheader">感谢使用我们的服务，仅差一步激活邮箱啦</span>
                    <table role="presentation" class="main">

                        <!-- START MAIN CONTENT AREA -->
                        {{if eq .TemplateType "email-verify"}}
                            {{template "email-verify-content" .}}
                        {{else if eq .TemplateType "password-reset"}}
                            {{template "password-reset-content" .}}
                        {{else}}
                            {{block "content" .}}{{end}}
                        {{end}}
                        <!-- END MAIN CONTENT AREA -->
                    </table>

                    <!-- START FOOTER -->
                    <!-- END FOOTER -->

                    <!-- END CENTERED WHITE CONTAINER -->
                </div>
            </td>
            <td>&nbsp;</td>
        </tr>
    </table>
    </body>

    </html>
{{end}}
