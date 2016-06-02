package downloader_test

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/golang/mock/gomock"
	dl "github.com/pivotal-cf/pcfdev-cli/downloader"
	"github.com/pivotal-cf/pcfdev-cli/downloader/mocks"
	"github.com/pivotal-cf/pcfdev-cli/pivnet"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Downloader", func() {
	var (
		downloader *dl.Downloader
		mockCtrl   *gomock.Controller
		mockClient *mocks.MockClient
		mockConfig *mocks.MockConfig
		mockFS     *mocks.MockFS
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockClient = mocks.NewMockClient(mockCtrl)
		mockConfig = mocks.NewMockConfig(mockCtrl)
		mockFS = mocks.NewMockFS(mockCtrl)

		downloader = &dl.Downloader{
			PivnetClient: mockClient,
			FS:           mockFS,
			Config:       mockConfig,
			ExpectedMD5:  "some-md5",
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("IsOVACurrent", func() {
		Context("when OVA does not exist", func() {
			It("should return false", func() {
				mockConfig.EXPECT().GetOVAPath().Return("some-ova-path", nil)
				mockFS.EXPECT().Exists("some-ova-path").Return(false, nil)

				current, err := downloader.IsOVACurrent()
				Expect(err).NotTo(HaveOccurred())
				Expect(current).To(BeFalse())

			})
		})
		Context("when OVA exists and has correct MD5", func() {
			It("should return true", func() {
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return("some-ova-path", nil),
					mockFS.EXPECT().Exists("some-ova-path").Return(true, nil),
					mockFS.EXPECT().MD5("some-ova-path").Return("some-md5", nil),
				)

				current, err := downloader.IsOVACurrent()
				Expect(err).NotTo(HaveOccurred())
				Expect(current).To(BeTrue())

			})
		})
		Context("when OVA exists and has incorrect MD5", func() {
			It("should return false", func() {
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return("some-ova-path", nil),
					mockFS.EXPECT().Exists("some-ova-path").Return(true, nil),
					mockFS.EXPECT().MD5("some-ova-path").Return("some-bad-md5", nil),
				)

				current, err := downloader.IsOVACurrent()
				Expect(err).NotTo(HaveOccurred())
				Expect(current).To(BeFalse())

			})
		})

		Context("when checking if the file exists fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return("some-ova-path", nil),
					mockFS.EXPECT().Exists("some-ova-path").Return(false, errors.New("some-error")),
				)

				_, err := downloader.IsOVACurrent()
				Expect(err).To(MatchError("some-error"))
			})
		})

		Context("when checking the MD5 fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return("some-ova-path", nil),
					mockFS.EXPECT().Exists("some-ova-path").Return(true, nil),
					mockFS.EXPECT().MD5("some-ova-path").Return("", errors.New("some-error")),
				)

				_, err := downloader.IsOVACurrent()
				Expect(err).To(MatchError("some-error"))
			})
		})

		Context("when getting the OVA path gives an error", func() {
			It("should return an error", func() {
				mockConfig.EXPECT().GetOVAPath().Return("some-ova-path", errors.New("some-error"))

				_, err := downloader.IsOVACurrent()
				Expect(err).To(MatchError("some-error"))
			})
		})
	})

	Describe("#Download", func() {
		Context("when file and partial file do not exist", func() {
			It("should download the file", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(false, nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-md5", nil),
					mockFS.EXPECT().Move(filepath.Join("some-path", "some-file.ova.partial"), filepath.Join("some-path", "some-file.ova")),
				)

				Expect(downloader.Download()).To(Succeed())
			})
		})

		Context("when partial file does exist", func() {
			It("should resume the download of the partial file", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(25), nil),
					mockClient.EXPECT().DownloadOVA(int64(25)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-md5", nil),
					mockFS.EXPECT().Move(filepath.Join("some-path", "some-file.ova.partial"), filepath.Join("some-path", "some-file.ova")),
				)

				Expect(downloader.Download()).To(Succeed())
			})
		})

		Context("when partial file is downloaded but the checksum is not valid and the re-download succeeds", func() {
			It("should move the file to the downloaded path", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(25), nil),
					mockClient.EXPECT().DownloadOVA(int64(25)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
					mockFS.EXPECT().RemoveFile(filepath.Join("some-path", "some-file.ova.partial")).Return(nil),

					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-md5", nil),
					mockFS.EXPECT().Move(filepath.Join("some-path", "some-file.ova.partial"), filepath.Join("some-path", "some-file.ova")),
				)

				Expect(downloader.Download()).To(Succeed())
			})
		})

		Context("when partial file is downloaded but the checksum is not valid and the re-download fails", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(25), nil),
					mockClient.EXPECT().DownloadOVA(int64(25)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
					mockFS.EXPECT().RemoveFile(filepath.Join("some-path", "some-file.ova.partial")).Return(nil),

					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
				)

				Expect(downloader.Download()).To(MatchError("download failed"))
			})
		})

		Context("when getting the OVA path returns an error", func() {
			It("should return an error", func() {
				mockConfig.EXPECT().GetOVAPath().Return("", errors.New("some-error"))

				Expect(downloader.Download()).To(MatchError("some-error"))
			})
		})

		Context("when creating the directory fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(errors.New("some-error")),
				)

				Expect(downloader.Download()).To(MatchError("some-error"))
			})
		})

		Context("when deleting files fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(errors.New("some-error")),
				)

				Expect(downloader.Download()).To(MatchError("some-error"))
			})
		})

		Context("when checking if the partial file exists", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(false, errors.New("some-error")),
				)

				Expect(downloader.Download()).To(MatchError("some-error"))
			})
		})

		Context("when checking the length of the partial file fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(0), errors.New("some-error")),
				)

				Expect(downloader.Download()).To(MatchError("some-error"))
			})
		})

		Context("when downloading the file fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(false, nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(nil, errors.New("some-error")),
				)

				Expect(downloader.Download()).To(MatchError("some-error"))
			})
		})

		Context("when saving the Pivnet API token fails", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(false, nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken().Return(errors.New("some-error")),
				)

				Expect(downloader.Download()).To(MatchError("some-error"))
			})
		})

		Context("when writing the downloaded file fails", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(false, nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(errors.New("some-error")),
				)

				Expect(downloader.Download()).To(MatchError("some-error"))
			})
		})

		Context("when checking the MD5 of the downloaded file fails", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(false, nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("", errors.New("some-error")),
				)

				Expect(downloader.Download()).To(MatchError("some-error"))
			})
		})

		Context("when removing the partial file fails", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(25), nil),
					mockClient.EXPECT().DownloadOVA(int64(25)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
					mockFS.EXPECT().RemoveFile(filepath.Join("some-path", "some-file.ova.partial")).Return(errors.New("some-error")),
				)

				Expect(downloader.Download()).To(MatchError("some-error"))
			})
		})

		Context("when the MD5 of a file download does not match the expected MD5", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(false, nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
				)

				Expect(downloader.Download()).To(MatchError("download failed"))
			})
		})

		Context("when downloading the file fails after downloading the partial file failed", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(25), nil),
					mockClient.EXPECT().DownloadOVA(int64(25)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
					mockFS.EXPECT().RemoveFile(filepath.Join("some-path", "some-file.ova.partial")).Return(nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(nil, errors.New("some-error")),
				)

				Expect(downloader.Download()).To(MatchError("some-error"))
			})
		})

		Context("when saving the Pivnet API token fails after downloading the partial file failed", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(25), nil),
					mockClient.EXPECT().DownloadOVA(int64(25)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
					mockFS.EXPECT().RemoveFile(filepath.Join("some-path", "some-file.ova.partial")).Return(nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken().Return(errors.New("some-error")),
				)

				Expect(downloader.Download()).To(MatchError("some-error"))
			})
		})

		Context("when writing the file fails after downloading the partial file failed", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(25), nil),
					mockClient.EXPECT().DownloadOVA(int64(25)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
					mockFS.EXPECT().RemoveFile(filepath.Join("some-path", "some-file.ova.partial")).Return(nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(errors.New("some-error")),
				)

				Expect(downloader.Download()).To(MatchError("some-error"))
			})
		})

		Context("when checking the MD5 of the file fails after downloading the partial file failed", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(25), nil),
					mockClient.EXPECT().DownloadOVA(int64(25)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
					mockFS.EXPECT().RemoveFile(filepath.Join("some-path", "some-file.ova.partial")).Return(nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("", errors.New("some-error")),
				)

				Expect(downloader.Download()).To(MatchError("some-error"))
			})
		})

		Context("when moving the file fails", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockConfig.EXPECT().GetOVAPath().Return(filepath.Join("some-path", "some-file.ova"), nil),
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().DeleteAllExcept("some-path", []string{"some-file.ova", "some-file.ova.partial"}).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(false, nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockConfig.EXPECT().SaveToken(),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-md5", nil),
					mockFS.EXPECT().Move(filepath.Join("some-path", "some-file.ova.partial"), filepath.Join("some-path", "some-file.ova")).Return(errors.New("some-error")),
				)

				Expect(downloader.Download()).To(MatchError("some-error"))
			})
		})
	})
})
