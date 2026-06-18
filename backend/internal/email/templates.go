package email

import "fmt"

func layout(preheader, bodyHTML, ctaURL, ctaLabel string) string {
	cta := ""
	if ctaURL != "" {
		cta = fmt.Sprintf(`
			<tr><td style="padding:28px 0 0 0;">
				<a href="%s" style="background:#E84142;color:#fff;text-decoration:none;font-weight:600;font-size:14px;padding:12px 22px;border-radius:10px;display:inline-block;">%s</a>
			</td></tr>`, ctaURL, ctaLabel)
	}
	return fmt.Sprintf(`<!doctype html>
<html><body style="margin:0;background:#0A0A0B;font-family:-apple-system,Segoe UI,Roboto,sans-serif;">
<span style="display:none;">%s</span>
<table width="100%%" cellpadding="0" cellspacing="0" style="background:#0A0A0B;padding:40px 0;">
<tr><td align="center">
<table width="480" cellpadding="0" cellspacing="0" style="background:#141416;border:1px solid #27272A;border-radius:16px;padding:32px;">
<tr><td style="color:#E84142;font-weight:700;font-size:13px;letter-spacing:.08em;text-transform:uppercase;">Team1 Blog</td></tr>
<tr><td style="padding-top:16px;color:#E4E4E7;font-size:15px;line-height:1.6;">%s</td></tr>
%s
</table>
</td></tr>
</table>
</body></html>`, preheader, bodyHTML, cta)
}

func InvitationEmail(role, registerURL string) (subject, html string) {
	subject = "You're invited to join Team1 Blog"
	body := fmt.Sprintf(`You've been invited to join the Team1 Blog Contributor Platform as a <strong>%s</strong>. This link expires in 72 hours.`, role)
	return subject, layout(subject, body, registerURL, "Accept Invitation")
}

func SubmissionAcknowledgedEmail(title, dashboardURL string) (subject, html string) {
	subject = fmt.Sprintf("Submitted for review: %s", title)
	body := fmt.Sprintf(`Your article <strong>%s</strong> was submitted for review. We'll notify you as soon as a moderator responds.`, title)
	return subject, layout(subject, body, dashboardURL, "View Status")
}

func ReviewFeedbackEmail(title, dashboardURL string) (subject, html string) {
	subject = fmt.Sprintf("Changes requested: %s", title)
	body := fmt.Sprintf(`A moderator left feedback on <strong>%s</strong>. Review the inline suggestions and resubmit when ready.`, title)
	return subject, layout(subject, body, dashboardURL, "View Feedback")
}

func ArticleApprovedEmail(title, dashboardURL string) (subject, html string) {
	subject = fmt.Sprintf("Editorially approved: %s", title)
	body := fmt.Sprintf(`Great news — <strong>%s</strong> has been editorially approved and is headed to graphic design.`, title)
	return subject, layout(subject, body, dashboardURL, "View Article")
}

func BannerReadyEmail(title, dashboardURL string) (subject, html string) {
	subject = fmt.Sprintf("Banner ready: %s", title)
	body := fmt.Sprintf(`The cover banner for <strong>%s</strong> is ready. It's now queued for publishing.`, title)
	return subject, layout(subject, body, dashboardURL, "View Article")
}

func ReadyToPublishEmail(title, dashboardURL string) (subject, html string) {
	subject = fmt.Sprintf("Ready to publish: %s", title)
	body := fmt.Sprintf(`<strong>%s</strong> is editorially approved with a banner attached — it's ready to go live.`, title)
	return subject, layout(subject, body, dashboardURL, "Publish Now")
}

func ArticlePublishedEmail(title, dashboardURL string) (subject, html string) {
	subject = fmt.Sprintf("Published: %s", title)
	body := fmt.Sprintf(`<strong>%s</strong> is now live on Team1 Blog. Payment release is in progress.`, title)
	return subject, layout(subject, body, dashboardURL, "View Article")
}

func PaymentInitiatedEmail(title, amount, dashboardURL string) (subject, html string) {
	subject = fmt.Sprintf("Payment initiated: %s", title)
	body := fmt.Sprintf(`A payment of <strong>$%s</strong> for <strong>%s</strong> has been sent to your Core wallet onchain. Confirmation usually takes a few moments.`, amount, title)
	return subject, layout(subject, body, dashboardURL, "View Payment")
}

func PaymentConfirmedEmail(title, amount, txHash, dashboardURL string) (subject, html string) {
	subject = fmt.Sprintf("Payment confirmed: %s", title)
	body := fmt.Sprintf(`Your payment of <strong>$%s</strong> for <strong>%s</strong> has been confirmed onchain.<br><br>Tx hash: <code>%s</code>`, amount, title, txHash)
	return subject, layout(subject, body, dashboardURL, "View Payment")
}

// PaymentFailedEmail goes to Super Admins, not the contributor - it means
// the onchain transfer didn't confirm and needs manual attention before
// retrying.
func PaymentFailedEmail(title, contributorName, txHash, dashboardURL string) (subject, html string) {
	subject = fmt.Sprintf("Payment failed to confirm: %s", title)
	body := fmt.Sprintf(
		`The payment for <strong>%s</strong> (contributor: %s) did not confirm onchain and needs review.<br><br>Tx hash: <code>%s</code>`,
		title, contributorName, txHash,
	)
	return subject, layout(subject, body, dashboardURL, "Review Payment")
}
