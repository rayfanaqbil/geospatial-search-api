package report

import (
	"context"
	"errors"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetDataLaporanMasukHariini(db *mongo.Database, waGroupId string) (msg string) {
	msg += "*Jumlah Laporan Hari ini:*\n"
	ranklist := GetRankDataLaporanHariini(db, TodayFilter(), waGroupId)
	for i, data := range ranklist {
		msg += strconv.Itoa(i+1) + ". " + data.Username + " : +" + strconv.Itoa(int(data.Poin)) + "\n"
	}

	return
}

func GenerateRekapMessageKemarinPerWAGroupID(db *mongo.Database, groupId string) (msg string, err error) {
	pushReportCounts, err := GetDataRepoMasukKemarinPerWaGroupID(db, groupId)
	if err != nil {
		return
	}
	laporanCounts, err := GetDataLaporanKemarinPerWAGroupID(db, groupId)
	if err != nil {
		return
	}
	mergedCounts := MergePhoneNumberCounts(pushReportCounts, laporanCounts)
	if len(mergedCounts) == 0 {
		err = errors.New("tidak ada aktifitas push dan laporan")
		return
	}
	msg = "*Laporan Penambahan Poin Total Kemarin :*\n"
	var phoneSlice []string
	for phoneNumber, info := range mergedCounts {
		msg += "✅ " + info.Name + " (" + phoneNumber + ") : +" + strconv.FormatFloat(info.Count, 'f', -1, 64) + "\n"
		if info.Count > 2 { //klo lebih dari 2 maka tidak akan dikurangi masuk ke daftra putih
			phoneSlice = append(phoneSlice, phoneNumber)
		}
	}

	if !HariLibur(GetDateKemarin()) { //kalo bukan kemaren hari libur maka akan ada pengurangan poin
		filter := bson.M{"wagroupid": groupId}
		var projectDocuments []model.Project
		projectDocuments, err = atdb.GetAllDoc[[]model.Project](db, "project", filter)
		if err != nil {
			return
		}
		msg += "\n*Laporan Pengurangan Poin Kemarin :*\n"

		// Buat map untuk menyimpan nomor telepon dari slice
		phoneMap := make(map[string]bool)

		// Masukkan semua nomor telepon dari slice ke dalam map
		for _, phoneNumber := range phoneSlice {
			phoneMap[phoneNumber] = true
		}
		// Buat map untuk melacak pengguna yang sudah diproses
		processedUsers := make(map[string]bool)

		// Iterasi melalui nomor telepon dalam dokumen MongoDB
		for _, doc := range projectDocuments {
			for _, member := range doc.Menu {
				phoneNumber := member.ID
				// Periksa apakah nomor telepon ada dalam map
				if _, exists := phoneMap[phoneNumber]; !exists {
					if !processedUsers[member.ID] {
						msg += "⛔ " + member.Name + " (" + member.ID + ") : -3\n"
						KurangPoinUserbyPhoneNumber(db, member.ID, 3)
						processedUsers[member.ID] = true
					}
				}
			}
		}
		msg += "\n\n*Klo pada hari kerja kurang dari 3 poin, maka dikurangi 3 poin ya ka. Cemunguddhh..*"
	} else {
		if HariLibur(GetDateSekarang()) {
			msg += "\n\n*Have a nice day :)*"
		} else {
			msg += "\n\n*Yuk bisa yuk... Semangat untuk hari ini...*"
		}

	}

	return
}

func GetDataRepoMasukKemarinPerWaGroupID(db *mongo.Database, groupId string) (phoneNumberCount map[string]PhoneNumberInfo, err error) {
	filter := bson.M{"_id": YesterdayFilter(), "project.wagroupid": groupId}
	pushrepodata, err := atdb.GetAllDoc[[]model.PushReport](db, "pushrepo", filter)
	if err != nil {
		return
	}
	phoneNumberCount = CountDuplicatePhoneNumbersWithName(pushrepodata)
	return
}

func GetDataLaporanKemarinPerWAGroupID(db *mongo.Database, waGroupId string) (phoneNumberCount map[string]PhoneNumberInfo, err error) {
	filter := bson.M{"_id": YesterdayFilter(), "project.wagroupid": waGroupId}
	laps, err := atdb.GetAllDoc[[]model.Laporan](db, "uxlaporan", filter)
	if err != nil {
		return
	}
	phoneNumberCount = CountDuplicatePhoneNumbersLaporan(laps)
	return
}

func GetRankDataLaporanHariini(db *mongo.Database, filterhari bson.M, waGroupId string) (ranklist []PushRank) {
	//uxlaporan := db.Collection("uxlaporan")
	// Create filter to query data for today
	filter := bson.M{"_id": filterhari, "project.wagroupid": waGroupId}
	//nopetugass, _ := atdb.GetAllDistinctDoc(db, filter, "nopetugas", "uxlaporan")
	laps, _ := atdb.GetAllDoc[[]model.Laporan](db, "uxlaporan", filter)
	print(len(laps))
	//ranklist := []PushRank{}
	for _, lap := range laps {
		if lap.Project.WAGroupID == waGroupId {
			ranklist = append(ranklist, PushRank{Username: lap.Petugas, Poin: 1})
		}
		//ranklist = append(ranklist, PushRank{Username: pushdata[0].Petugas, Poin: float64(len(pushdata))})

	}
	return
}

func GetDataLaporanMasukHarian(db *mongo.Database) (msg string) {
	msg += "*Jumlah Laporan Hari Ini :*\n"
	ranklist := GetRankDataLayananHarian(db, TodayFilter())
	for i, data := range ranklist {
		msg += strconv.Itoa(i+1) + ". " + data.Username + " : " + strconv.Itoa(data.TotalCommit) + "\n"
	}

	return
}
func GetRankDataLayananHarian(db *mongo.Database, filterhari bson.M) (ranklist []PushRank) {
	pushrepo := db.Collection("uxlaporan")
	// Create filter to query data for today
	filter := bson.M{"_id": filterhari}
	usernamelist, _ := atdb.GetAllDistinctDoc(db, filter, "petugas", "uxlaporan")
	//ranklist := []PushRank{}
	for _, username := range usernamelist {
		filter := bson.M{"petugas": username, "_id": filterhari}
		// Query the database
		var pushdata []model.Laporan
		cur, err := pushrepo.Find(context.Background(), filter)
		if err != nil {
			return
		}
		if err = cur.All(context.Background(), &pushdata); err != nil {
			return
		}
		defer cur.Close(context.Background())
		if len(pushdata) > 0 {
			ranklist = append(ranklist, PushRank{Username: username.(string), TotalCommit: len(pushdata)})
		}
	}
	sort.SliceStable(ranklist, func(i, j int) bool {
		return ranklist[i].TotalCommit > ranklist[j].TotalCommit
	})
	return
}

func GetDataRepoMasukKemarinBukanLibur(db *mongo.Database) (msg string) {
	msg += "*Laporan Jumlah Push Repo Hari Ini :*\n"
	pushrepo := db.Collection("pushrepo")
	// Create filter to query data for today
	filter := bson.M{"_id": YesterdayNotLiburFilter()}
	usernamelist, _ := atdb.GetAllDistinctDoc(db, filter, "username", "pushrepo")
	for _, username := range usernamelist {
		filter := bson.M{"username": username, "_id": YesterdayNotLiburFilter()}
		// Query the database
		var pushdata []model.PushReport
		cur, err := pushrepo.Find(context.Background(), filter)
		if err != nil {
			return
		}
		if err = cur.All(context.Background(), &pushdata); err != nil {
			return
		}
		defer cur.Close(context.Background())
		if len(pushdata) > 0 {
			msg += "*" + username.(string) + " : " + strconv.Itoa(len(pushdata)) + "*\n"
			for j, push := range pushdata {
				msg += strconv.Itoa(j+1) + ". " + strings.TrimSpace(push.Message) + "\n"

			}
		}
	}
	return
}

func GetDataRepoMasukHariIni(db *mongo.Database, groupId string) (msg string) {
	msg += "*Laporan Penambahan Poin dari Jumlah Push Repo Hari ini :*\n"
	pushrepo := db.Collection("pushrepo")
	// Create filter to query data for today
	filter := bson.M{"_id": TodayFilter(), "project.wagroupid": groupId}
	usernamelist, _ := atdb.GetAllDistinctDoc(db, filter, "username", "pushrepo")
	for _, username := range usernamelist {
		filter := bson.M{"username": username, "_id": TodayFilter()}
		// Query the database
		var pushdata []model.PushReport
		cur, err := pushrepo.Find(context.Background(), filter)
		if err != nil {
			return
		}
		if err = cur.All(context.Background(), &pushdata); err != nil {
			return
		}
		defer cur.Close(context.Background())
		if len(pushdata) > 0 {
			msg += "*" + username.(string) + " : +" + strconv.Itoa(len(pushdata)) + "*\n"
			for j, push := range pushdata {
				msg += strconv.Itoa(j+1) + ". " + strings.TrimSpace(push.Message) + "\n"

			}
		}
	}
	return
}

func GetDataRepoMasukHarian(db *mongo.Database) (msg string) {
	msg += "*Laporan Jumlah Push Repo Hari Ini :*\n"
	pushrepo := db.Collection("pushrepo")
	// Create filter to query data for today
	filter := bson.M{"_id": TodayFilter()}
	usernamelist, _ := atdb.GetAllDistinctDoc(db, filter, "username", "pushrepo")
	for _, username := range usernamelist {
		filter := bson.M{"username": username, "_id": TodayFilter()}
		// Query the database
		var pushdata []model.PushReport
		cur, err := pushrepo.Find(context.Background(), filter)
		if err != nil {
			return
		}
		if err = cur.All(context.Background(), &pushdata); err != nil {
			return
		}
		defer cur.Close(context.Background())
		if len(pushdata) > 0 {
			msg += "*" + username.(string) + " : " + strconv.Itoa(len(pushdata)) + "*\n"
			for j, push := range pushdata {
				msg += strconv.Itoa(j+1) + ". " + strings.TrimSpace(push.Message) + "\n"

			}
		}
	}
	return
}

func GetRankDataRepoMasukHarian(db *mongo.Database, filterhari bson.M) (ranklist []PushRank) {
	pushrepo := db.Collection("pushrepo")
	// Create filter to query data for today
	filter := bson.M{"_id": filterhari}
	usernamelist, _ := atdb.GetAllDistinctDoc(db, filter, "username", "pushrepo")
	//ranklist := []PushRank{}
	for _, username := range usernamelist {
		filter := bson.M{"username": username, "_id": filterhari}
		cur, err := pushrepo.Find(context.Background(), filter)
		if err != nil {
			log.Println("Failed to find pushrepo data:", err)
			return
		}

		defer cur.Close(context.Background())

		repoCommits := make(map[string]int)
		for cur.Next(context.Background()) {
			var report model.PushReport
			if err := cur.Decode(&report); err != nil {
				log.Println("Failed to decode pushrepo data:", err)
				return
			}
			repoCommits[report.Repo]++
		}

		if len(repoCommits) > 0 {
			totalCommits := 0
			for _, count := range repoCommits {
				totalCommits += count
			}
			ranklist = append(ranklist, PushRank{Username: username.(string), TotalCommit: totalCommits, Repos: repoCommits})
		}
	}
	sort.SliceStable(ranklist, func(i, j int) bool {
		return ranklist[i].TotalCommit > ranklist[j].TotalCommit
	})
	return
}

func GetDateSekarang() (datesekarang time.Time) {
	// Definisi lokasi waktu sekarang
	location, _ := time.LoadLocation("Asia/Jakarta")

	t := time.Now().In(location) //.Truncate(24 * time.Hour)
	datesekarang = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	return
}

func TodayFilter() bson.M {
	return bson.M{
		"$gte": primitive.NewObjectIDFromTimestamp(GetDateSekarang()),
		"$lt":  primitive.NewObjectIDFromTimestamp(GetDateSekarang().Add(24 * time.Hour)),
	}
}

func YesterdayNotLiburFilter() bson.M {
	return bson.M{
		"$gte": primitive.NewObjectIDFromTimestamp(GetDateKemarinBukanHariLibur()),
		"$lt":  primitive.NewObjectIDFromTimestamp(GetDateKemarinBukanHariLibur().Add(24 * time.Hour)),
	}
}

func YesterdayFilter() bson.M {
	return bson.M{
		"$gte": primitive.NewObjectIDFromTimestamp(GetDateKemarin()),
		"$lt":  primitive.NewObjectIDFromTimestamp(GetDateKemarin().Add(24 * time.Hour)),
	}
}

func GetDateKemarinBukanHariLibur() (datekemarinbukanlibur time.Time) {
	// Definisi lokasi waktu sekarang
	location, _ := time.LoadLocation("Asia/Jakarta")
	n := -1
	t := time.Now().AddDate(0, 0, n).In(location) //.Truncate(24 * time.Hour)
	for HariLibur(t) {
		n -= 1
		t = time.Now().AddDate(0, 0, n).In(location)
	}

	datekemarinbukanlibur = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	return
}

func GetDateKemarin() (datekemarin time.Time) {
	// Definisi lokasi waktu sekarang
	location, _ := time.LoadLocation("Asia/Jakarta")
	n := -1
	t := time.Now().AddDate(0, 0, n).In(location) //.Truncate(24 * time.Hour)
	datekemarin = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	return
}

func HariLibur(thedate time.Time) (libur bool) {
	wekkday := thedate.Weekday()
	inhari := int(wekkday)
	if inhari == 0 || inhari == 6 {
		libur = true
	}
	tglskr := thedate.Format("2006-01-02")
	tgl := int(thedate.Month())
	urltarget := "https://dayoffapi.vercel.app/api?month=" + strconv.Itoa(tgl)
	_, hasil, _ := atapi.Get[[]NewLiburNasional](urltarget)
	for _, v := range hasil {
		if v.Tanggal == tglskr {
			libur = true
		}
	}
	return
}

func Last3DaysFilter() bson.M {
	tigaHariLalu := GetDateSekarang().Add(-72 * time.Hour) // 3 * 24 hours
	now := GetDateSekarang()
	return bson.M{
		"$gte": primitive.NewObjectIDFromTimestamp(tigaHariLalu),
		"$lt":  primitive.NewObjectIDFromTimestamp(now),
	}
}
