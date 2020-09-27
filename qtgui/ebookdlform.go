package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/visualfc/goqt/ui"

	edl "github.com/sndnvaps/ebookdownloader"
	ebook "github.com/sndnvaps/ebookdownloader/ebook-sources"
)

type EbookdlForm struct {
	m           *ui.QWidget
	sfNameLabel *ui.QLabel
	bookIdLabel *ui.QLabel
	proxyLabel  *ui.QLabel

	bookIdInput *ui.QLineEdit
	proxyInput  *ui.QLineEdit

	websiteComboBox *ui.QComboBox //用于设置默认的下载网站

	outputTypeLayout   *ui.QGridLayout
	outputEpubCheckBok *ui.QCheckBox
	outputMobiCheckBok *ui.QCheckBox
	outputTxtCheckBok  *ui.QCheckBox

	downloadProgressBar *ui.QProgressBar

	downloadButton *ui.QPushButton
}

func NewEbookDlForm() (*EbookdlForm, error) {
	w := &EbookdlForm{}
	w.m = ui.NewWidget()
	w.m.SetFixedSizeWithWidthHeight(400, 400) //设置固定大小的窗口

	file := ui.NewFileWithName(":/forms/ebookdlform.ui")

	if !file.Open(ui.QIODevice_ReadOnly) {
		return nil, errors.New("error load ui")
	}

	loader := ui.NewUiLoader()
	formWidget := loader.Load(file)
	if formWidget == nil {
		return nil, errors.New("error load form widget")
	}

	w.sfNameLabel = ui.NewLabelFromDriver(formWidget.FindChild("softwarename"))
	w.bookIdLabel = ui.NewLabelFromDriver(formWidget.FindChild("bookid"))
	w.proxyLabel = ui.NewLabelFromDriver(formWidget.FindChild("proxy"))

	w.bookIdInput = ui.NewLineEditFromDriver(formWidget.FindChild("bookidInput"))
	w.proxyInput = ui.NewLineEditFromDriver(formWidget.FindChild("proxyInput"))

	w.websiteComboBox = ui.NewComboBoxFromDriver(formWidget.FindChild("defWebsiteCB"))

	websiteLists := []string{"xsbiquge.com", "biduo.cc", "xixiwx.com", "booktxt.net", "biquwu.cc", "999xs.com", "23us.la"}
	w.websiteComboBox.AddItems(websiteLists)

	w.outputTypeLayout = ui.NewGridLayoutFromDriver(formWidget.FindChild("OutputTypeLayout"))

	//设置选项
	w.outputTxtCheckBok = ui.NewCheckBoxFromDriver(formWidget.FindChild("OutputTxt"))
	w.outputMobiCheckBok = ui.NewCheckBoxFromDriver(formWidget.FindChild("OutputMobi"))
	w.outputEpubCheckBok = ui.NewCheckBoxFromDriver(formWidget.FindChild("OutputEpub"))

	//w.outputTxtCheckBok.IsChecked() //判断是否被打勾

	w.downloadProgressBar = ui.NewProgressBarFromDriver(formWidget.FindChild("progressBar"))

	var bookinfo edl.BookInfo              //初始化变量
	var EBDLInterface edl.EBookDLInterface //初始化接口

	//var metainfo edl.Meta //用于保存小说的meta信息

	w.downloadButton = ui.NewPushButtonFromDriver(formWidget.FindChild("StartButton"))
	w.downloadButton.OnClicked(func() {
		//设置下载进度条
		/*
			w.downloadProgressBar.SetRange(0, 1000)
			for i := 1; i < 1000+1; {
				i++
				w.downloadProgressBar.SetValue((int32)(i))
			}
		*/
		bookid := ""
		//proxy := ""
		if strings.Compare(w.bookIdInput.Text(), "") != 0 {
			bookid = w.bookIdInput.Text()
		}
		/*
			if strings.Compare(w.proxyInput.Text(), "") != 0 {
				proxy = w.proxyInput.Text()
			}
		*/
		if w.websiteComboBox.CurrentText() == "xsbiquge.com" {
			xsbiquge := ebook.NewXSBiquge()
			EBDLInterface = xsbiquge //实例化接口
		} else if w.websiteComboBox.CurrentText() == "biduo.cc" {
			biduo := ebook.NewBiDuo()
			EBDLInterface = biduo //实例化接口
		} else if w.websiteComboBox.CurrentText() == "booktxt.net" {
			booktxt := ebook.NewBookTXT()
			EBDLInterface = booktxt //实例化接口
		} else {
			messagebox := ui.NewMessageBox()
			messagebox.SetText("必须要选择一个下载源！")
			messagebox.Show()
		}

		w.downloadProgressBar.SetRange(0, 100)
		bookinfo = EBDLInterface.GetBookInfo(bookid, "")

		w.downloadProgressBar.SetValue(1)
		bookinfo = EBDLInterface.DownloadChapters(bookinfo, "") //下载小说章节内容
		w.downloadProgressBar.SetValue(25)

		if w.outputTxtCheckBok.IsChecked() {
			bookinfo.GenerateTxt()
			w.downloadProgressBar.SetValue(50)
		}
		if w.outputMobiCheckBok.IsChecked() {
			bookinfo.SetKindleEbookType(true /* isMobi */, false /* isAzw3 */)
			bookinfo.GenerateMobi()
			w.downloadProgressBar.SetValue(60)
		}

		if w.outputEpubCheckBok.IsChecked() {
			bookinfo.GenerateEPUB()
			w.downloadProgressBar.SetValue(70)
		}
		w.downloadProgressBar.SetValue(100)
		messagebox := ui.NewMessageBox()
		outputInfo := fmt.Sprintf("小说名：%s\n作者：%s\n简介：\n\t%s", bookinfo.Name, bookinfo.Author, bookinfo.Description)
		messagebox.SetText(outputInfo)
		messagebox.Show()
		w.downloadProgressBar.Reset()
	})

	layout := ui.NewVBoxLayout()
	layout.AddWidget(formWidget)
	w.m.SetLayout(layout)

	w.m.SetWindowTitle("Ebookdownloader")
	return w, nil
}