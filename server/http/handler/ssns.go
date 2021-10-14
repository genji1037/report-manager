package handler

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"os"
	"report-manager/alg"
	"report-manager/db"
	"report-manager/proxy"
	"report-manager/server/http/respond"
	"strconv"
	"strings"
	"time"
)

type GetSSNSReportResp struct {
	Dat         string          `json:"dat"`
	Seq         int             `json:"seq"`
	BonusAmount decimal.Decimal `json:"bonus_amount"`
	LinkAmount  decimal.Decimal `json:"link_amount"`
}

func GetSSNSReport(c *gin.Context) {
	date := c.Query("date")
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		respond.BadRequest(c, http.StatusBadRequest, err.Error())
		return
	}
	seqStr := c.Query("seq")
	seq, err := strconv.Atoi(seqStr)
	if err != nil {
		respond.BadRequest(c, http.StatusBadRequest, fmt.Sprintf("bad seq %s", seqStr))
		return
	}

	var cacheNotFound bool
	r := db.SSNSReport{Dat: date, Seq: seq}
	if err := r.GetByDateSeq(); err != nil {
		if err == gorm.ErrRecordNotFound {
			cacheNotFound = true
		} else {
			respond.InternalError(c, err)
			return
		}
	}

	if cacheNotFound {
		var datePrefix string
		if seq == 1 {
			yesterdayDate, _ := alg.DateAdd(date, -1)
			datePrefix = strings.ReplaceAll(yesterdayDate, "-", "") + "2"
		} else if seq == 2 {
			datePrefix = strings.ReplaceAll(date, "-", "") + "0"
		} else {
			respond.BadRequest(c, http.StatusBadRequest, fmt.Sprintf("bad seq %s", seqStr))
			return
		}

		assertErr := func(err error) bool {
			if err != nil {
				if err == proxy.ErrFileNotFound {
					respond.BadRequest(c, http.StatusBadRequest, fmt.Sprintf("network reward not found at %s_%d", date, seq))
					return true
				}
				respond.InternalError(c, err)
				return true
			}
			return false
		}

		networkRewardSum, err := getRewardFileSummary("network_" + datePrefix)
		if assertErr(err) {
			return
		}

		linkRewardSum, err := getRewardFileSummary("link_" + datePrefix)
		if assertErr(err) {
			return
		}

		r.BonusAmount = networkRewardSum.Add(linkRewardSum)

		var linkFilePrefix string
		if seq == 1 {
			linkFilePrefix = fmt.Sprintf("relation_%s+%d", date, 4)
		} else {
			linkFilePrefix = fmt.Sprintf("relation_%s+%d", date, 16)
		}
		linkFileBs, err := proxy.GetSSNSFile("server-secret-social-network/file/consensus/relation", linkFilePrefix)
		if assertErr(err) {
			return
		}
		doubledLinkSum, err := getLinkSummary(linkFileBs)
		if err != nil {
			respond.InternalError(c, err)
			return
		}
		r.LinkAmount = doubledLinkSum.Div(decimal.NewFromFloat(2))
		r.Create()
	}

	respond.Success(c, GetSSNSReportResp{
		Dat:         r.Dat,
		Seq:         r.Seq,
		BonusAmount: r.BonusAmount,
		LinkAmount:  r.LinkAmount,
	})
}

func getLinkSummary(bs []byte) (decimal.Decimal, error) {
	rd := bufio.NewReader(bytes.NewBuffer(bs))
	var sum decimal.Decimal
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return decimal.Decimal{}, fmt.Errorf("read line failed: %v", err)
			}
		}
		arr := strings.Split(string(line), ",")
		if len(arr) < 3 {
			return decimal.Decimal{}, fmt.Errorf("bad line %s", string(line))
		}
		d, err := decimal.NewFromString(arr[2])
		if err != nil {
			return decimal.Decimal{}, fmt.Errorf("bad line %s", string(line))
		}
		sum = sum.Add(d)
	}
	return sum, nil
}

// getRewardFile returns temp dir
func getRewardFileSummary(prefix string) (decimal.Decimal, error) {
	// reward file
	reward, err := proxy.GetSSNSFile("server-secret-social-network/file/consensus/reward", prefix)
	if err != nil {
		return decimal.Decimal{}, err
	}

	tempDir := os.TempDir()
	if !strings.HasSuffix(tempDir, "/") {
		tempDir += "/"
	}
	tempPath := tempDir + prefix
	f, err := os.Create(tempPath)
	if err != nil {
		return decimal.Decimal{}, err
	}
	_, err = f.Write(reward)
	if err != nil {
		return decimal.Decimal{}, err
	}
	f.Close()
	return solveZippedRewardFile(tempPath)
}

// returns summary of rewards.
func solveZippedRewardFile(filePath string) (decimal.Decimal, error) {
	rc, err := zip.OpenReader(filePath)
	if err != nil {
		return decimal.Decimal{}, err
	}
	defer rc.Close()

	if len(rc.Reader.File) == 0 || rc.Reader.File[0] == nil {
		return decimal.Decimal{}, errors.New("empty zip file")
	}
	f, err := rc.Reader.File[0].Open()
	if err != nil {
		return decimal.Decimal{}, err
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	var sum decimal.Decimal
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return decimal.Decimal{}, fmt.Errorf("read line failed: %v", err)
			}
		}
		arr := strings.Split(string(line), ",")
		if len(arr) < 2 {
			return decimal.Decimal{}, fmt.Errorf("bad line %s", string(line))
		}
		d, err := decimal.NewFromString(arr[1])
		if err != nil {
			return decimal.Decimal{}, fmt.Errorf("bad line %s", string(line))
		}
		sum = sum.Add(d)
	}
	return sum, nil
}
