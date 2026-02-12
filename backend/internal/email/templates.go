package email

import (
	"fmt"
	"unicode"
)

func titleCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func wrapHTML(title, content string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><title>%s</title></head>
<body style="margin:0;padding:0;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#f4f4f5;">
<div style="max-width:560px;margin:40px auto;background:#fff;border-radius:12px;overflow:hidden;">
<div style="background:#18181b;padding:24px 32px;">
<h1 style="margin:0;color:#fff;font-size:20px;font-weight:700;">Creatrid</h1>
</div>
<div style="padding:32px;">
%s
</div>
<div style="padding:16px 32px;background:#fafafa;text-align:center;">
<p style="margin:0;font-size:12px;color:#a1a1aa;">Creatrid — Creator Passport</p>
</div>
</div>
</body>
</html>`, title, content)
}

func WelcomeEmail(name, username, profileURL string) (subject, body string) {
	subject = "Welcome to Creatrid!"
	content := fmt.Sprintf(`
<h2 style="margin:0 0 16px;color:#18181b;font-size:22px;">Welcome, %s!</h2>
<p style="color:#52525b;line-height:1.6;">Your Creator Passport is ready. You can now connect your social accounts, build your Creator Score, and share your verified profile with brands and collaborators.</p>
<div style="margin:24px 0;">
<a href="%s" style="display:inline-block;background:#18181b;color:#fff;padding:12px 24px;border-radius:8px;text-decoration:none;font-weight:600;font-size:14px;">View Your Profile</a>
</div>
<p style="color:#52525b;line-height:1.6;">Here's what to do next:</p>
<ol style="color:#52525b;line-height:1.8;padding-left:20px;">
<li>Connect your social platforms (YouTube, GitHub, Twitter, etc.)</li>
<li>Customize your profile theme and add custom links</li>
<li>Share your profile link with brands and collaborators</li>
</ol>
<p style="color:#a1a1aa;font-size:13px;margin-top:24px;">Your profile: <a href="%s" style="color:#18181b;">%s</a></p>
`, name, profileURL, profileURL, profileURL)
	body = wrapHTML(subject, content)
	return
}

func ConnectionAlertEmail(name, platform, platformUsername string) (subject, body string) {
	subject = fmt.Sprintf("%s connected to Creatrid", titleCase(platform))
	content := fmt.Sprintf(`
<h2 style="margin:0 0 16px;color:#18181b;font-size:22px;">New Connection: %s</h2>
<p style="color:#52525b;line-height:1.6;">Hi %s, your <strong>%s</strong> account (<strong>%s</strong>) has been successfully connected to your Creator Passport.</p>
<p style="color:#52525b;line-height:1.6;">Your Creator Score has been recalculated to reflect this new connection. Connect more platforms to boost your score further.</p>
<div style="margin:24px 0;">
<a href="https://creatrid.com/connections" style="display:inline-block;background:#18181b;color:#fff;padding:12px 24px;border-radius:8px;text-decoration:none;font-weight:600;font-size:14px;">View Connections</a>
</div>
`, titleCase(platform), name, titleCase(platform), platformUsername)
	body = wrapHTML(subject, content)
	return
}

func AccountDeletedEmail(name string) (subject, body string) {
	subject = "Your Creatrid account has been deleted"
	content := fmt.Sprintf(`
<h2 style="margin:0 0 16px;color:#18181b;font-size:22px;">Account Deleted</h2>
<p style="color:#52525b;line-height:1.6;">Hi %s, your Creatrid account and all associated data have been permanently deleted. This includes:</p>
<ul style="color:#52525b;line-height:1.8;padding-left:20px;">
<li>Your profile and Creator Score</li>
<li>All connected social accounts</li>
<li>Analytics data (profile views and link clicks)</li>
<li>Collaboration requests</li>
</ul>
<p style="color:#52525b;line-height:1.6;">If you didn't request this, please contact us immediately at <a href="mailto:support@creatrid.com" style="color:#18181b;">support@creatrid.com</a>.</p>
<p style="color:#a1a1aa;font-size:13px;margin-top:24px;">You can always create a new account at <a href="https://creatrid.com" style="color:#18181b;">creatrid.com</a>.</p>
`, name)
	body = wrapHTML(subject, content)
	return
}

func EmailVerificationEmail(name, verifyURL string) (subject, body string) {
	subject = "Verify your email on Creatrid"
	content := fmt.Sprintf(`
<h2 style="margin:0 0 16px;color:#18181b;font-size:22px;">Verify Your Email</h2>
<p style="color:#52525b;line-height:1.6;">Hi %s, please verify your email address to complete your Creator Passport and earn 10 bonus points on your Creator Score.</p>
<div style="margin:24px 0;">
<a href="%s" style="display:inline-block;background:#18181b;color:#fff;padding:12px 24px;border-radius:8px;text-decoration:none;font-weight:600;font-size:14px;">Verify Email</a>
</div>
<p style="color:#a1a1aa;font-size:13px;margin-top:24px;">This link expires in 24 hours. If you didn't request this, you can safely ignore this email.</p>
`, name, verifyURL)
	body = wrapHTML(subject, content)
	return
}

func WeeklyDigestEmail(name string, totalViews, viewsThisWeek, totalClicks, connectionCount int, score *int) (subject, body string) {
	subject = "Your weekly Creatrid digest"
	scoreStr := "—"
	if score != nil {
		scoreStr = fmt.Sprintf("%d/100", *score)
	}
	content := fmt.Sprintf(`
<h2 style="margin:0 0 16px;color:#18181b;font-size:22px;">Weekly Digest</h2>
<p style="color:#52525b;line-height:1.6;">Hi %s, here's your weekly summary:</p>
<table style="width:100%%;border-collapse:collapse;margin:20px 0;">
<tr>
<td style="padding:12px 16px;background:#fafafa;border-radius:8px 0 0 0;text-align:center;">
<p style="margin:0;font-size:24px;font-weight:700;color:#18181b;">%d</p>
<p style="margin:4px 0 0;font-size:12px;color:#a1a1aa;">Views This Week</p>
</td>
<td style="padding:12px 16px;background:#fafafa;text-align:center;">
<p style="margin:0;font-size:24px;font-weight:700;color:#18181b;">%d</p>
<p style="margin:4px 0 0;font-size:12px;color:#a1a1aa;">Total Views</p>
</td>
<td style="padding:12px 16px;background:#fafafa;border-radius:0 8px 0 0;text-align:center;">
<p style="margin:0;font-size:24px;font-weight:700;color:#18181b;">%d</p>
<p style="margin:4px 0 0;font-size:12px;color:#a1a1aa;">Link Clicks</p>
</td>
</tr>
<tr>
<td style="padding:12px 16px;background:#fafafa;border-radius:0 0 0 8px;text-align:center;">
<p style="margin:0;font-size:24px;font-weight:700;color:#18181b;">%d</p>
<p style="margin:4px 0 0;font-size:12px;color:#a1a1aa;">Connections</p>
</td>
<td colspan="2" style="padding:12px 16px;background:#fafafa;border-radius:0 0 8px 0;text-align:center;">
<p style="margin:0;font-size:24px;font-weight:700;color:#18181b;">%s</p>
<p style="margin:4px 0 0;font-size:12px;color:#a1a1aa;">Creator Score</p>
</td>
</tr>
</table>
<div style="margin:24px 0;">
<a href="https://creatrid.com/dashboard" style="display:inline-block;background:#18181b;color:#fff;padding:12px 24px;border-radius:8px;text-decoration:none;font-weight:600;font-size:14px;">View Dashboard</a>
</div>
`, name, viewsThisWeek, totalViews, totalClicks, connectionCount, scoreStr)
	body = wrapHTML(subject, content)
	return
}
