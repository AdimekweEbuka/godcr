package modal

import (
	"fmt"
	"image/color"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const Info = "info_modal"

type InfoModal struct {
	*load.Load
	randomID        string
	modal           decredmaterial.Modal
	keyEvent        chan *key.Event
	enterKeyPressed bool

	dialogIcon *decredmaterial.Icon

	dialogTitle    string
	subtitle       string
	customTemplate []layout.Widget

	positiveButtonText    string
	positiveButtonClicked func()
	btnPositve            decredmaterial.Button

	negativeButtonText    string
	negativeButtonClicked func()
	btnNegative           decredmaterial.Button

	checkbox decredmaterial.CheckBoxStyle

	textInputEditor decredmaterial.Editor
	showEditor      bool
	isEnabled       bool

	titleAlignment, btnAlignment layout.Direction

	isCancelable bool
}

func NewInfoModal(l *load.Load) *InfoModal {
	in := &InfoModal{
		Load:         l,
		randomID:     fmt.Sprintf("%s-%d", Info, decredmaterial.GenerateRandomNumber()),
		modal:        *l.Theme.ModalFloatTitle(),
		btnPositve:   l.Theme.OutlineButton("Yes"),
		btnNegative:  l.Theme.OutlineButton("No"),
		keyEvent:     l.Receiver.KeyEvents,
		isCancelable: true,
		btnAlignment: layout.E,
	}

	in.btnPositve.Font.Weight = text.Medium
	in.btnNegative.Font.Weight = text.Medium

	in.textInputEditor = l.Theme.Editor(new(widget.Editor), "")
	in.textInputEditor.Editor.SingleLine, in.textInputEditor.Editor.Submit = true, true

	return in
}

func (in *InfoModal) ModalID() string {
	return in.randomID
}

func (in *InfoModal) Show() {
	in.ShowModal(in)
}

func (in *InfoModal) Dismiss() {
	in.DismissModal(in)
}

func (in *InfoModal) OnResume() {
	in.Load.EnableKeyEvent = true
}

func (in *InfoModal) OnDismiss() {
	in.Load.EnableKeyEvent = false
}

func (in *InfoModal) ShowTextEditor(showTextEditor bool) *InfoModal {
	in.showEditor = showTextEditor
	return in
}

func (in *InfoModal) TextInputEditor(editor decredmaterial.Editor) *InfoModal {
	in.textInputEditor = editor
	return in
}

func (in *InfoModal) SetCancelable(min bool) *InfoModal {
	in.isCancelable = min
	return in
}

func (in *InfoModal) SetContentAlignment(title, btn layout.Direction) *InfoModal {
	in.titleAlignment = title
	in.btnAlignment = btn
	return in
}

func (in *InfoModal) Icon(icon *decredmaterial.Icon) *InfoModal {
	in.dialogIcon = icon
	return in
}

func (in *InfoModal) CheckBox(checkbox decredmaterial.CheckBoxStyle) *InfoModal {
	in.checkbox = checkbox
	return in
}

func (in *InfoModal) Title(title string) *InfoModal {
	in.dialogTitle = title
	return in
}

func (in *InfoModal) Body(subtitle string) *InfoModal {
	in.subtitle = subtitle
	return in
}

func (in *InfoModal) PositiveButton(text string, clicked func()) *InfoModal {
	in.positiveButtonText = text
	in.positiveButtonClicked = clicked
	return in
}

func (in *InfoModal) PositiveButtonStyle(background, text color.NRGBA) *InfoModal {
	in.btnPositve.Background, in.btnPositve.Color = background, text
	return in
}

func (in *InfoModal) NegativeButton(text string, clicked func()) *InfoModal {
	in.negativeButtonText = text
	in.negativeButtonClicked = clicked
	return in
}

// for backwards compatibilty
func (in *InfoModal) SetupWithTemplate(template string) *InfoModal {
	title := in.dialogTitle
	subtitle := in.subtitle
	var customTemplate []layout.Widget
	switch template {
	case TransactionDetailsInfoTemplate:
		title = "How to copy"
		customTemplate = transactionDetailsInfo(in.Theme)
	case SignMessageInfoTemplate:
		customTemplate = signMessageInfo(in.Theme)
	case VerifyMessageInfoTemplate:
		customTemplate = verifyMessageInfo(in.Theme)
	case PrivacyInfoTemplate:
		title = "How to use the mixer?"
		customTemplate = privacyInfo(in.Theme)
	case SetupMixerInfoTemplate:
		customTemplate = setupMixerInfo(in.Theme)
	case WalletBackupInfoTemplate:
		customTemplate = backupInfo(in.Theme)
	case AllowUnmixedSpendingTemplate:
		title = "Confirm to allow spending from unmixed accounts"
		customTemplate = allowUnspendUnmixedAcct(in.Theme)
	}

	in.dialogTitle = title
	in.subtitle = subtitle
	in.customTemplate = customTemplate
	return in
}

func (in *InfoModal) handleEnterKeypress() {
	select {
	case event := <-in.keyEvent:
		if (event.Name == key.NameReturn || event.Name == key.NameEnter) && event.State == key.Press && in.customTemplate != nil {
			in.enterKeyPressed = true
		}
	default:
	}
}

func (in *InfoModal) Handle() {
	if in.showEditor {
		if editorsNotEmpty(in.textInputEditor.Editor) {
			in.btnPositve.Background = in.Theme.Color.Danger
			in.isEnabled = true
		} else {
			in.btnPositve.Background = in.Theme.Color.InactiveGray
			in.isEnabled = false
		}
	}

	for in.btnPositve.Clicked() {
		if !in.showEditor {
			in.DismissModal(in)
			in.positiveButtonClicked()
		} else {
			in.positiveButtonClicked()
		}
	}

	for in.btnNegative.Clicked() {
		if !in.showEditor && !editorsNotEmpty(in.textInputEditor.Editor) {
			in.DismissModal(in)
			in.negativeButtonClicked()
		} else {
			in.negativeButtonClicked()
		}
	}

	if in.modal.BackdropClicked(in.isCancelable) {
		in.Dismiss()
	}

	if in.checkbox.CheckBox != nil {
		in.btnNegative.SetEnabled(in.checkbox.CheckBox.Value)
	}
}

func (in *InfoModal) Layout(gtx layout.Context) D {
	icon := func(gtx C) D {
		if in.dialogIcon == nil {
			return layout.Dimensions{}
		}

		return layout.Inset{Top: values.MarginPadding10, Bottom: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return in.dialogIcon.Layout(gtx, values.MarginPadding50)
			})
		})
	}

	checkbox := func(gtx C) D {
		if in.checkbox.CheckBox == nil {
			return layout.Dimensions{}
		}

		return layout.Inset{Top: values.MarginPaddingMinus5, Left: values.MarginPaddingMinus5}.Layout(gtx, func(gtx C) D {
			in.checkbox.TextSize = values.TextSize14
			in.checkbox.Color = in.Theme.Color.GrayText1
			in.checkbox.IconColor = in.Theme.Color.Gray2
			if in.checkbox.CheckBox.Value {
				in.checkbox.IconColor = in.Theme.Color.Primary
			}
			return in.checkbox.Layout(gtx)
		})
	}

	textEditor := func(gtx C) D {
		if !in.showEditor {
			return layout.Dimensions{}
		}

		return layout.Inset{Top: values.MarginPaddingMinus5, Left: values.MarginPaddingMinus5}.Layout(gtx, func(gtx C) D {
			in.textInputEditor.TextSize = values.TextSize14
			in.textInputEditor.EditorStyle.Color = in.Theme.Color.Black
			in.textInputEditor.EditorStyle.Font.Weight = text.Weight(in.Theme.TextSize.U)
			return in.textInputEditor.Layout(gtx)
		})
	}

	subtitle := func(gtx C) D {
		text := in.Theme.Body1(in.subtitle)
		text.Color = in.Theme.Color.GrayText2
		return text.Layout(gtx)
	}

	var w []layout.Widget

	// Every section of the dialog is optional
	if in.dialogIcon != nil {
		w = append(w, icon)
	}

	if in.dialogTitle != "" {
		w = append(w, in.titleLayout())
	}

	if in.subtitle != "" {
		w = append(w, subtitle)
	}

	if in.customTemplate != nil {
		w = append(w, in.customTemplate...)
	}

	if in.showEditor {
		w = append(w, textEditor)
	}

	if in.checkbox.CheckBox != nil {
		w = append(w, checkbox)
	}

	if in.negativeButtonText != "" || in.positiveButtonText != "" {
		w = append(w, in.actionButtonsLayout())
	}

	return in.modal.Layout(gtx, w)
}

func (in *InfoModal) titleLayout() layout.Widget {
	return func(gtx C) D {
		t := in.Theme.H6(in.dialogTitle)
		t.Font.Weight = text.SemiBold
		return in.titleAlignment.Layout(gtx, t.Layout)
	}
}

func (in *InfoModal) actionButtonsLayout() layout.Widget {
	return func(gtx C) D {
		return in.btnAlignment.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if in.negativeButtonText == "" {
						return layout.Dimensions{}
					}

					in.btnNegative.Text = in.negativeButtonText
					return layout.Inset{Right: values.MarginPadding5}.Layout(gtx, in.btnNegative.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					if in.positiveButtonText == "" {
						return layout.Dimensions{}
					}

					in.btnPositve.Text = in.positiveButtonText
					return in.btnPositve.Layout(gtx)
				}),
			)
		})
	}
}
