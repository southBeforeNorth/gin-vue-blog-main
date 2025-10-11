{{define "password-reset-content"}}
    <tr>
        <td class="wrapper">
            <table role="presentation" border="0" cellpadding="0" cellspacing="0">
                <tr>
                    <td>
                        <p>👋&nbsp; 您好~ {{.UserName}} ~ </p>
                        <p>🔐&nbsp; 您正在申请修改密码，请使用以下验证码完成密码重置。</p>
                        <p>📬&nbsp; 您的验证码是：<strong style="font-size: 24px; color: #007bff;">{{.Code}}</strong></p>
                        <p>⏰&nbsp; 验证码有效期为5分钟，请及时使用。</p>
                        <p>⚠️&nbsp; 如果您没有申请修改密码，请忽略此邮件。</p>
                    </td>
                </tr>
            </table>
        </td>
    </tr>
{{end}}
