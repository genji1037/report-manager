package radar_otc

type RealName struct{}

func (RealName) CountWaitingRealNames() (int, error) {
	var result struct {
		Cnt int
	}
	sql := `
SELECT 
    count(0) as cnt
FROM
    deposits d
        LEFT JOIN
    real_names r ON d.uid = r.uid AND d.state = 'paid'
WHERE
    r.state = 'wait'
`
	err := gormDb.Raw(sql).Scan(&result).Error
	return result.Cnt, err
}
