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
	return wrapHTMLWithPreheader(title, "", content)
}

func wrapHTMLWithPreheader(title, preheader, content string) string {
	preheaderHTML := ""
	if preheader != "" {
		preheaderHTML = fmt.Sprintf(`<div style="display:none;font-size:1px;color:#f4f4f5;line-height:1px;max-height:0;max-width:0;opacity:0;overflow:hidden;">%s</div>`, preheader)
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><title>%s</title></head>
<body style="margin:0;padding:0;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#f4f4f5;">
%s
<div style="max-width:560px;margin:40px auto;background:#fff;border-radius:12px;overflow:hidden;">
<div style="background:#18181b;padding:24px 32px;">
<h1 style="margin:0;color:#fff;font-size:20px;font-weight:700;">Creatrid</h1>
</div>
<div style="padding:32px;">
%s
</div>
<div style="padding:16px 32px;background:#fafafa;text-align:center;">
<p style="margin:0 0 8px;font-size:12px;color:#a1a1aa;">Creatrid &mdash; Creator Passport</p>
<p style="margin:0;font-size:11px;color:#d4d4d8;">You received this email because you have a Creatrid account. To manage your email preferences, visit your <a href="https://creatrid.com/settings" style="color:#a1a1aa;text-decoration:underline;">account settings</a>.</p>
</div>
</div>
</body>
</html>`, title, preheaderHTML, content)
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

func CollaborationRequestEmail(recipientName, senderName, senderUsername, message, inboxURL string) (subject, body string) {
	subject = fmt.Sprintf("You have a new collaboration request from %s", senderName)
	msgHTML := ""
	if message != "" {
		msgHTML = fmt.Sprintf(`
<div style="margin:16px 0;padding:16px;background:#fafafa;border-radius:8px;border-left:4px solid #18181b;">
<p style="margin:0;color:#52525b;line-height:1.6;font-style:italic;">"%s"</p>
</div>`, message)
	}
	content := fmt.Sprintf(`
<h2 style="margin:0 0 16px;color:#18181b;font-size:22px;">New Collaboration Request</h2>
<p style="color:#52525b;line-height:1.6;">Hi %s, <strong>%s</strong> (@%s) wants to collaborate with you!</p>
%s
<div style="margin:24px 0;">
<a href="%s" style="display:inline-block;background:#18181b;color:#fff;padding:12px 24px;border-radius:8px;text-decoration:none;font-weight:600;font-size:14px;">View Request</a>
</div>
<p style="color:#a1a1aa;font-size:13px;margin-top:24px;">You can accept or decline this request from your collaboration inbox.</p>
`, recipientName, senderName, senderUsername, msgHTML, inboxURL)
	body = wrapHTMLWithPreheader(subject, fmt.Sprintf("%s wants to collaborate with you on Creatrid", senderName), content)
	return
}

func CollaborationResponseEmail(senderName, recipientName, status, collabURL string) (subject, body string) {
	subject = fmt.Sprintf("%s %s your collaboration request", recipientName, status)
	color := "#ef4444"
	if status == "accepted" {
		color = "#22c55e"
	}
	content := fmt.Sprintf(`
<h2 style="margin:0 0 16px;color:#18181b;font-size:22px;">Collaboration %s</h2>
<p style="color:#52525b;line-height:1.6;">Hi %s, <strong>%s</strong> has <span style="color:%s;font-weight:600;">%s</span> your collaboration request.</p>
<div style="margin:24px 0;">
<a href="%s" style="display:inline-block;background:#18181b;color:#fff;padding:12px 24px;border-radius:8px;text-decoration:none;font-weight:600;font-size:14px;">View Details</a>
</div>
`, titleCase(status), senderName, recipientName, color, status, collabURL)
	body = wrapHTMLWithPreheader(subject, fmt.Sprintf("%s %s your collaboration request", recipientName, status), content)
	return
}

func ContentFlaggedEmail(creatorName, contentTitle, reason, vaultURL string) (subject, body string) {
	subject = fmt.Sprintf("Your content '%s' has been flagged", contentTitle)
	content := fmt.Sprintf(`
<h2 style="margin:0 0 16px;color:#18181b;font-size:22px;">Content Flagged</h2>
<p style="color:#52525b;line-height:1.6;">Hi %s, your content <strong>%s</strong> has been flagged for review.</p>
<div style="margin:16px 0;padding:16px;background:#fef2f2;border-radius:8px;border-left:4px solid #ef4444;">
<p style="margin:0;color:#991b1b;line-height:1.6;"><strong>Reason:</strong> %s</p>
</div>
<p style="color:#52525b;line-height:1.6;">Our team will review the flagged content. If it violates our community guidelines, it may be removed. You can update your content metadata to address the issue.</p>
<div style="margin:24px 0;">
<a href="%s" style="display:inline-block;background:#18181b;color:#fff;padding:12px 24px;border-radius:8px;text-decoration:none;font-weight:600;font-size:14px;">Go to Content Vault</a>
</div>
`, creatorName, contentTitle, reason, vaultURL)
	body = wrapHTMLWithPreheader(subject, fmt.Sprintf("Your content '%s' needs attention", contentTitle), content)
	return
}

func LicensePurchasedEmail(sellerName, buyerName, contentTitle, licenseType string, priceCents int, earningsURL string) (subject, body string) {
	subject = "Someone purchased a license for your content"
	priceStr := fmt.Sprintf("$%.2f", float64(priceCents)/100)
	content := fmt.Sprintf(`
<h2 style="margin:0 0 16px;color:#18181b;font-size:22px;">License Purchased!</h2>
<p style="color:#52525b;line-height:1.6;">Hi %s, great news! <strong>%s</strong> just purchased a <strong>%s</strong> license for your content:</p>
<div style="margin:16px 0;padding:16px;background:#f0fdf4;border-radius:8px;border-left:4px solid #22c55e;">
<p style="margin:0 0 8px;color:#18181b;font-weight:600;">%s</p>
<p style="margin:0;color:#52525b;">License type: %s</p>
<p style="margin:4px 0 0;color:#22c55e;font-weight:700;font-size:18px;">%s</p>
</div>
<div style="margin:24px 0;">
<a href="%s" style="display:inline-block;background:#18181b;color:#fff;padding:12px 24px;border-radius:8px;text-decoration:none;font-weight:600;font-size:14px;">View Earnings</a>
</div>
`, sellerName, buyerName, licenseType, contentTitle, licenseType, priceStr, earningsURL)
	body = wrapHTMLWithPreheader(subject, fmt.Sprintf("%s purchased a %s license for '%s'", buyerName, licenseType, contentTitle), content)
	return
}

func ProfileMilestoneEmail(name, milestone, dashboardURL string) (subject, body string) {
	subject = fmt.Sprintf("Congratulations! You've reached %s", milestone)
	content := fmt.Sprintf(`
<h2 style="margin:0 0 16px;color:#18181b;font-size:22px;">Milestone Reached!</h2>
<p style="color:#52525b;line-height:1.6;">Hi %s, congratulations on reaching a new milestone:</p>
<div style="margin:16px 0;padding:24px;background:linear-gradient(135deg,#fafafa,#f0f0f0);border-radius:12px;text-align:center;">
<p style="margin:0;font-size:28px;font-weight:700;color:#18181b;">%s</p>
</div>
<p style="color:#52525b;line-height:1.6;">Keep building your Creator Passport to unlock even more achievements. Connect more platforms, engage with your audience, and grow your Creator Score.</p>
<div style="margin:24px 0;">
<a href="%s" style="display:inline-block;background:#18181b;color:#fff;padding:12px 24px;border-radius:8px;text-decoration:none;font-weight:600;font-size:14px;">View Dashboard</a>
</div>
`, name, milestone, dashboardURL)
	body = wrapHTMLWithPreheader(subject, fmt.Sprintf("You've reached %s on Creatrid!", milestone), content)
	return
}
