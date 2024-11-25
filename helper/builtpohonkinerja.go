package helper

import (
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/opdmaster"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"sort"
)

func BuildTematikResponse(pohonMap map[int]map[int][]domain.PohonKinerja, tematik domain.PohonKinerja) pohonkinerja.TematikResponse {
	tematikResp := pohonkinerja.TematikResponse{
		Id:         tematik.Id,
		Parent:     nil,
		Tema:       tematik.NamaPohon,
		JenisPohon: tematik.JenisPohon,
		LevelPohon: tematik.LevelPohon,
		Keterangan: tematik.Keterangan,
		Indikators: ConvertToIndikatorResponses(tematik.Indikator),
	}

	var childs []interface{}

	// Tambahkan strategic (level 4) yang memiliki parent level 0
	if strategics := pohonMap[4][tematik.Id]; len(strategics) > 0 {
		for _, strategic := range strategics {
			strategicResp := BuildStrategicResponse(pohonMap, strategic)
			childs = append(childs, strategicResp)
		}
	}

	// Tambahkan subtematik (level 1)
	if subTematiks := pohonMap[1][tematik.Id]; len(subTematiks) > 0 {
		for _, subTematik := range subTematiks {
			subTematikResp := BuildSubTematikResponse(pohonMap, subTematik)
			childs = append(childs, subTematikResp)
		}
	}

	tematikResp.Child = childs
	return tematikResp
}

func BuildSubTematikResponse(pohonMap map[int]map[int][]domain.PohonKinerja, subTematik domain.PohonKinerja) pohonkinerja.SubtematikResponse {
	subTematikResp := pohonkinerja.SubtematikResponse{
		Id:         subTematik.Id,
		Parent:     subTematik.Parent,
		Tema:       subTematik.NamaPohon,
		JenisPohon: subTematik.JenisPohon,
		LevelPohon: subTematik.LevelPohon,
		Keterangan: subTematik.Keterangan,
		Indikators: ConvertToIndikatorResponses(subTematik.Indikator),
	}

	var childs []interface{}

	// Tambahkan strategic (level 4) yang memiliki parent level 1
	if strategics := pohonMap[4][subTematik.Id]; len(strategics) > 0 {
		for _, strategic := range strategics {
			strategicResp := BuildStrategicResponse(pohonMap, strategic)
			childs = append(childs, strategicResp)
		}
	}

	// Tambahkan subsubtematik (level 2)
	if subSubTematiks := pohonMap[2][subTematik.Id]; len(subSubTematiks) > 0 {
		for _, subSubTematik := range subSubTematiks {
			subSubTematikResp := BuildSubSubTematikResponse(pohonMap, subSubTematik)
			childs = append(childs, subSubTematikResp)
		}
	}

	subTematikResp.Child = childs
	return subTematikResp
}

func BuildSubSubTematikResponse(pohonMap map[int]map[int][]domain.PohonKinerja, subSubTematik domain.PohonKinerja) pohonkinerja.SubSubTematikResponse {
	subSubTematikResp := pohonkinerja.SubSubTematikResponse{
		Id:         subSubTematik.Id,
		Parent:     subSubTematik.Parent,
		Tema:       subSubTematik.NamaPohon,
		JenisPohon: subSubTematik.JenisPohon,
		LevelPohon: subSubTematik.LevelPohon,
		Keterangan: subSubTematik.Keterangan,
		Indikators: ConvertToIndikatorResponses(subSubTematik.Indikator),
	}

	var childs []interface{}

	// Tambahkan strategic (level 4) yang memiliki parent level 2
	if strategics := pohonMap[4][subSubTematik.Id]; len(strategics) > 0 {
		for _, strategic := range strategics {
			strategicResp := BuildStrategicResponse(pohonMap, strategic)
			childs = append(childs, strategicResp)
		}
	}

	// Tambahkan supersubtematik (level 3)
	if superSubTematiks := pohonMap[3][subSubTematik.Id]; len(superSubTematiks) > 0 {
		for _, superSubTematik := range superSubTematiks {
			superSubTematikResp := BuildSuperSubTematikResponse(pohonMap, superSubTematik)
			childs = append(childs, superSubTematikResp)
		}
	}

	subSubTematikResp.Child = childs
	return subSubTematikResp
}

func BuildSuperSubTematikResponse(pohonMap map[int]map[int][]domain.PohonKinerja, superSubTematik domain.PohonKinerja) pohonkinerja.SuperSubTematikResponse {
	superSubTematikResp := pohonkinerja.SuperSubTematikResponse{
		Id:         superSubTematik.Id,
		Parent:     superSubTematik.Parent,
		Tema:       superSubTematik.NamaPohon,
		JenisPohon: superSubTematik.JenisPohon,
		LevelPohon: superSubTematik.LevelPohon,
		Keterangan: superSubTematik.Keterangan,
		Indikators: ConvertToIndikatorResponses(superSubTematik.Indikator),
	}

	var childs []interface{}

	// Tambahkan strategic (level 4) yang memiliki parent level 3
	if strategics := pohonMap[4][superSubTematik.Id]; len(strategics) > 0 {
		for _, strategic := range strategics {
			strategicResp := BuildStrategicResponse(pohonMap, strategic)
			childs = append(childs, strategicResp)
		}
	}

	superSubTematikResp.Childs = childs
	return superSubTematikResp
}

func BuildStrategicResponse(pohonMap map[int]map[int][]domain.PohonKinerja, strategic domain.PohonKinerja) pohonkinerja.StrategicResponse {
	strategicResp := pohonkinerja.StrategicResponse{
		Id:         strategic.Id,
		Parent:     strategic.Parent,
		Strategi:   strategic.NamaPohon,
		JenisPohon: strategic.JenisPohon,
		LevelPohon: strategic.LevelPohon,
		Keterangan: strategic.Keterangan,
		Indikators: ConvertToIndikatorResponses(strategic.Indikator),
		KodeOpd: &opdmaster.OpdResponseForAll{
			KodeOpd: strategic.KodeOpd,
			NamaOpd: strategic.NamaOpd,
		},
	}

	var childs []interface{}

	// Tambahkan tactical (level 5) ke childs
	if tacticals := pohonMap[5][strategic.Id]; len(tacticals) > 0 {
		// Urutkan tactical berdasarkan Id
		sort.Slice(tacticals, func(i, j int) bool {
			return tacticals[i].Id < tacticals[j].Id
		})

		for _, tactical := range tacticals {
			tacticalResp := BuildTacticalResponse(pohonMap, tactical)
			childs = append(childs, tacticalResp)
		}
	}

	strategicResp.Childs = childs
	return strategicResp
}

func BuildTacticalResponse(pohonMap map[int]map[int][]domain.PohonKinerja, tactical domain.PohonKinerja) pohonkinerja.TacticalResponse {
	tacticalResp := pohonkinerja.TacticalResponse{
		Id:         tactical.Id,
		Parent:     tactical.Parent,
		Strategi:   tactical.NamaPohon,
		JenisPohon: tactical.JenisPohon,
		LevelPohon: tactical.LevelPohon,
		Keterangan: &tactical.Keterangan,
		Indikators: ConvertToIndikatorResponses(tactical.Indikator),
	}

	// Tambahkan data OPD jika ada
	if tactical.KodeOpd != "" {
		tacticalResp.KodeOpd = &opdmaster.OpdResponseForAll{
			KodeOpd: tactical.KodeOpd,
			NamaOpd: tactical.NamaOpd,
		}
	}

	var childs []interface{}

	// Tambahkan operational ke childs
	if operationals := pohonMap[6][tactical.Id]; len(operationals) > 0 {
		for _, operational := range operationals {
			operationalResp := BuildOperationalResponse(pohonMap, operational)
			childs = append(childs, operationalResp)
		}
	}

	tacticalResp.Childs = childs
	return tacticalResp
}

func BuildOperationalResponse(pohonMap map[int]map[int][]domain.PohonKinerja, operational domain.PohonKinerja) pohonkinerja.OperationalResponse {
	operationalResp := pohonkinerja.OperationalResponse{
		Id:         operational.Id,
		Parent:     operational.Parent,
		Strategi:   operational.NamaPohon,
		JenisPohon: operational.JenisPohon,
		LevelPohon: operational.LevelPohon,
		Keterangan: &operational.Keterangan,
		Indikators: ConvertToIndikatorResponses(operational.Indikator),
	}

	// Tambahkan data OPD jika ada
	if operational.KodeOpd != "" {
		operationalResp.KodeOpd = &opdmaster.OpdResponseForAll{
			KodeOpd: operational.KodeOpd,
			NamaOpd: operational.NamaOpd,
		}
	}

	var childs []interface{}

	// Cek level berikutnya (operational-n)
	nextLevel := operational.LevelPohon + 1
	if operationalNs := pohonMap[nextLevel][operational.Id]; len(operationalNs) > 0 {
		// Urutkan berdasarkan Id
		sort.Slice(operationalNs, func(i, j int) bool {
			return operationalNs[i].Id < operationalNs[j].Id
		})

		for _, opN := range operationalNs {
			operationalNResp := BuildOperationalNResponse(pohonMap, opN)
			childs = append(childs, operationalNResp)
		}
	}

	operationalResp.Childs = childs
	return operationalResp
}

func BuildOperationalNResponse(pohonMap map[int]map[int][]domain.PohonKinerja, operationalN domain.PohonKinerja) pohonkinerja.OperationalNResponse {
	operationalNResp := pohonkinerja.OperationalNResponse{
		Id:         operationalN.Id,
		Parent:     operationalN.Parent,
		Strategi:   operationalN.NamaPohon,
		JenisPohon: operationalN.JenisPohon,
		LevelPohon: operationalN.LevelPohon,
		Keterangan: &operationalN.Keterangan,
		Indikators: ConvertToIndikatorResponses(operationalN.Indikator),
	}

	// Tambahkan data OPD jika ada
	if operationalN.KodeOpd != "" {
		operationalNResp.KodeOpd = &opdmaster.OpdResponseForAll{
			KodeOpd: operationalN.KodeOpd,
			NamaOpd: operationalN.NamaOpd,
		}
	}

	// Cek level berikutnya secara rekursif
	nextLevel := operationalN.LevelPohon + 1
	if nextOperationalNs := pohonMap[nextLevel][operationalN.Id]; len(nextOperationalNs) > 0 {
		// Urutkan berdasarkan Id
		sort.Slice(nextOperationalNs, func(i, j int) bool {
			return nextOperationalNs[i].Id < nextOperationalNs[j].Id
		})

		var childs []pohonkinerja.OperationalNResponse
		for _, nextOpN := range nextOperationalNs {
			childResp := BuildOperationalNResponse(pohonMap, nextOpN)
			childs = append(childs, childResp)
		}
		operationalNResp.Childs = childs
	}

	return operationalNResp
}

func BuildSubTematikResponseLimited(pohonMap map[int]map[int][]domain.PohonKinerja, subTematik domain.PohonKinerja) pohonkinerja.SubtematikResponse {
	subTematikResp := pohonkinerja.SubtematikResponse{
		Id:         subTematik.Id,
		Parent:     subTematik.Parent,
		Tema:       subTematik.NamaPohon,
		JenisPohon: subTematik.JenisPohon,
		LevelPohon: subTematik.LevelPohon,
		Keterangan: subTematik.Keterangan,
		Indikators: ConvertToIndikatorResponses(subTematik.Indikator),
	}

	var childs []interface{}

	// Hanya tambahkan strategic (level 4) yang memiliki parent level 1
	if strategics := pohonMap[4][subTematik.Id]; len(strategics) > 0 {
		for _, strategic := range strategics {
			strategicResp := BuildStrategicResponse(pohonMap, strategic)
			childs = append(childs, strategicResp)
		}
	}

	subTematikResp.Child = childs
	return subTematikResp
}
