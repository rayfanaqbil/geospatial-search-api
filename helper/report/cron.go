package report

import (
	"encoding/base64"
	"errors"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/whatsauth"
	"github.com/gocroot/model"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RekapMeetingKemarin(db *mongo.Database) (err error) {
	filter := bson.M{"_id": YesterdayFilter()}
	wagroupidlist, err := atdb.GetAllDistinctDoc(db, filter, "project.wagroupid", "uxlaporan")
	if err != nil {
		return
	}
	if len(wagroupidlist) == 0 {
		return
	}
	for _, gid := range wagroupidlist { //iterasi di setiap wa group
		// Type assertion to convert any to string
		groupID, ok := gid.(string)
		if !ok {
			err = errors.New("wagroupid is not a string, skipping this iteration")
			continue
		}
		filter := bson.M{"wagroupid": groupID}
		var projectDocuments []model.Project
		projectDocuments, err = atdb.GetAllDoc[[]model.Project](db, "project", filter)
		if err != nil {
			continue
		}
		for _, project := range projectDocuments {
			var base64pdf, md string
			base64pdf, md, err = GetPDFandMDMeeting(db, project.Name)
			if err != nil {
				continue
			}
			dt := &itmodel.DocumentMessage{
				To:        groupID,
				IsGroup:   true,
				Base64Doc: base64pdf,
				Filename:  project.Name + ".pdf",
				Caption:   "Berikut ini rekap rapat kemaren ya kak untuk project " + project.Name,
			}
			_, _, err = atapi.PostStructWithToken[model.Response]("Token", config.WAAPIToken, dt, config.WAAPIDocMessage)
			if err != nil {
				continue
			}
			//upload file markdown ke log repo untuk tipe rapat
			if project.RepoLogName != "" {
				// Encode string ke base64
				encodedString := base64.StdEncoding.EncodeToString([]byte(md))

				// Format markdown dengan base64 string
				//markdownContent := fmt.Sprintf("```base64\n%s\n```", encodedString)
				dt := model.LogInfo{
					PhoneNumber: project.Owner.PhoneNumber,
					Alias:       project.Owner.Name,
					FileName:    "README.md",
					RepoOrg:     project.RepoOrg,
					RepoName:    project.RepoLogName,
					Base64Str:   encodedString,
				}
				var conf model.Config
				conf, err = atdb.GetOneDoc[model.Config](db, "config", bson.M{"phonenumber": "62895601060000"})
				if err != nil {
					continue
				}

				//masalahnya disini pake token pribadi. kalo user awangga tidak masuk ke repo maka ga bisa
				atapi.PostStructWithToken[model.LogInfo]("secret", conf.LeaflySecret, dt, conf.LeaflyURL)
			}
		}
	}

	return

}

func RekapPagiHari(db *mongo.Database) (err error) {
	filter := bson.M{"_id": YesterdayFilter()}
	wagroupidlist, err := atdb.GetAllDistinctDoc(db, filter, "project.wagroupid", "pushrepo")
	if err != nil {
		return errors.New("Gagal Query Distinct project.wagroupid: " + err.Error())
	}

	var lastErr error // Variabel untuk menyimpan kesalahan terakhir

	for _, gid := range wagroupidlist { // Iterasi di setiap wa group
		// Type assertion to convert any to string
		groupID, ok := gid.(string)
		if !ok {
			lastErr = errors.New("wagroupid is not a string")
			continue
		}
		var msg string
		msg, err = GenerateRekapMessageKemarinPerWAGroupID(db, groupID)
		if err != nil {
			lastErr = errors.New("Gagal Membuat Rekapitulasi perhitungan per wa group id: " + err.Error())
			continue
		}
		dt := &whatsauth.TextMessage{
			To:       groupID,
			IsGroup:  true,
			Messages: msg,
		}
		var resp model.Response
		_, resp, err = atapi.PostStructWithToken[model.Response]("Token", config.WAAPIToken, dt, config.WAAPIMessage)
		if err != nil {
			lastErr = errors.New("Tidak berhak: " + err.Error() + ", " + resp.Info)
			continue
		}
	}

	if lastErr != nil {
		return lastErr
	}

	return nil
}
