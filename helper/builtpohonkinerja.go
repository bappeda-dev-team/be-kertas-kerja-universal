package helper

import (
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/opdmaster"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
)

func BuildTematikResponse(pohonMap map[int]map[int][]domain.PohonKinerja, tematik domain.PohonKinerja) pohonkinerja.TematikResponse {
	tematikResp := pohonkinerja.TematikResponse{
		Id:         tematik.Id,
		Parent:     nil,
		Tema:       tematik.NamaPohon,
		Keterangan: tematik.Keterangan,
		Indikators: ConvertToIndikatorResponses(tematik.Indikator),
	}

	if subTematiks := pohonMap[1][tematik.Id]; len(subTematiks) > 0 {
		var subTematikResponses []pohonkinerja.SubtematikResponse
		for _, subTematik := range subTematiks {
			subTematikResp := BuildSubTematikResponse(pohonMap, subTematik)
			subTematikResponses = append(subTematikResponses, subTematikResp)
		}
		tematikResp.SubTematiks = subTematikResponses
	} else {
		if strategics := pohonMap[4][tematik.Id]; len(strategics) > 0 {
			tematikResp.Strategics = BuildStrategicResponses(pohonMap, strategics)
		}
	}

	return tematikResp
}

func BuildSubTematikResponse(pohonMap map[int]map[int][]domain.PohonKinerja, subTematik domain.PohonKinerja) pohonkinerja.SubtematikResponse {
	subTematikResp := pohonkinerja.SubtematikResponse{
		Id:         subTematik.Id,
		Parent:     subTematik.Parent,
		Tema:       subTematik.NamaPohon,
		Keterangan: subTematik.Keterangan,
		Indikators: ConvertToIndikatorResponses(subTematik.Indikator),
	}

	if subSubTematiks := pohonMap[2][subTematik.Id]; len(subSubTematiks) > 0 {
		var subSubTematikResponses []pohonkinerja.SubSubTematikResponse
		for _, subSubTematik := range subSubTematiks {
			subSubTematikResp := BuildSubSubTematikResponse(pohonMap, subSubTematik)
			subSubTematikResponses = append(subSubTematikResponses, subSubTematikResp)
		}
		subTematikResp.SubSubTematiks = subSubTematikResponses
	} else {
		if strategics := pohonMap[4][subTematik.Id]; len(strategics) > 0 {
			subTematikResp.Strategics = BuildStrategicResponses(pohonMap, strategics)
		}
	}

	return subTematikResp
}

func BuildSubSubTematikResponse(pohonMap map[int]map[int][]domain.PohonKinerja, subSubTematik domain.PohonKinerja) pohonkinerja.SubSubTematikResponse {
	subSubTematikResp := pohonkinerja.SubSubTematikResponse{
		Id:         subSubTematik.Id,
		Parent:     subSubTematik.Parent,
		Tema:       subSubTematik.NamaPohon,
		Keterangan: subSubTematik.Keterangan,
		Indikators: ConvertToIndikatorResponses(subSubTematik.Indikator),
	}

	if superSubTematiks := pohonMap[3][subSubTematik.Id]; len(superSubTematiks) > 0 {
		var superSubTematikResponses []pohonkinerja.SuperSubTematikResponse
		for _, superSubTematik := range superSubTematiks {
			superSubTematikResp := BuildSuperSubTematikResponse(pohonMap, superSubTematik)
			superSubTematikResponses = append(superSubTematikResponses, superSubTematikResp)
		}
		subSubTematikResp.SuperSubTematiks = superSubTematikResponses
	} else {
		if strategics := pohonMap[4][subSubTematik.Id]; len(strategics) > 0 {
			subSubTematikResp.Strategics = BuildStrategicResponses(pohonMap, strategics)
		}
	}

	return subSubTematikResp
}

func BuildSuperSubTematikResponse(pohonMap map[int]map[int][]domain.PohonKinerja, superSubTematik domain.PohonKinerja) pohonkinerja.SuperSubTematikResponse {
	superSubTematikResp := pohonkinerja.SuperSubTematikResponse{
		Id:         superSubTematik.Id,
		Parent:     superSubTematik.Parent,
		Tema:       superSubTematik.NamaPohon,
		Keterangan: superSubTematik.Keterangan,
		Indikators: ConvertToIndikatorResponses(superSubTematik.Indikator),
	}

	if strategics := pohonMap[4][superSubTematik.Id]; len(strategics) > 0 {
		superSubTematikResp.Strategics = BuildStrategicResponses(pohonMap, strategics)
	}

	return superSubTematikResp
}

func BuildStrategicResponses(pohonMap map[int]map[int][]domain.PohonKinerja, strategics []domain.PohonKinerja) []pohonkinerja.StrategicResponse {
	var responses []pohonkinerja.StrategicResponse
	for _, strategic := range strategics {
		strategicResp := pohonkinerja.StrategicResponse{
			Id:         strategic.Id,
			Parent:     strategic.Parent,
			Strategi:   strategic.NamaPohon,
			Keterangan: strategic.Keterangan,
			KodeOpd: opdmaster.OpdResponseForAll{
				KodeOpd: strategic.KodeOpd,
				NamaOpd: strategic.NamaOpd, // Pastikan field ini terisi
			},
			Indikators: ConvertToIndikatorResponses(strategic.Indikator),
		}

		if tacticals := pohonMap[5][strategic.Id]; len(tacticals) > 0 {
			strategicResp.Tacticals = BuildTacticalResponses(pohonMap, tacticals)
		}

		responses = append(responses, strategicResp)
	}
	return responses
}

func BuildTacticalResponses(pohonMap map[int]map[int][]domain.PohonKinerja, tacticals []domain.PohonKinerja) []pohonkinerja.TacticalResponse {
	var responses []pohonkinerja.TacticalResponse
	for _, tactical := range tacticals {
		keterangan := &tactical.Keterangan
		tacticalResp := pohonkinerja.TacticalResponse{
			Id:         tactical.Id,
			Parent:     tactical.Parent,
			Strategi:   tactical.NamaPohon,
			Keterangan: keterangan,
			KodeOpd: opdmaster.OpdResponseForAll{
				KodeOpd: tactical.KodeOpd,
				NamaOpd: tactical.NamaOpd, // Pastikan field ini terisi
			},
			Indikators: ConvertToIndikatorResponses(tactical.Indikator),
		}

		if operationals := pohonMap[6][tactical.Id]; len(operationals) > 0 {
			tacticalResp.Operationals = BuildOperationalResponses(operationals)
		}

		responses = append(responses, tacticalResp)
	}
	return responses
}

func BuildOperationalResponses(operationals []domain.PohonKinerja) []pohonkinerja.OperationalResponse {
	var responses []pohonkinerja.OperationalResponse
	for _, operational := range operationals {
		keterangan := &operational.Keterangan
		operationalResp := pohonkinerja.OperationalResponse{
			Id:         operational.Id,
			Parent:     operational.Parent,
			Strategi:   operational.NamaPohon,
			Keterangan: keterangan,
			KodeOpd: opdmaster.OpdResponseForAll{
				KodeOpd: operational.KodeOpd,
				NamaOpd: operational.NamaOpd, // Pastikan field ini terisi
			},
			Indikators: ConvertToIndikatorResponses(operational.Indikator),
		}
		responses = append(responses, operationalResp)
	}
	return responses
}

// BuildTematikResponseLimited hanya membangun response untuk level 0 dan 1
func BuildTematikResponseLimited(pohonMap map[int]map[int][]domain.PohonKinerja, tematik domain.PohonKinerja) pohonkinerja.TematikResponse {
	tematikResp := pohonkinerja.TematikResponse{
		Id:          tematik.Id,
		Parent:      nil,
		Tema:        tematik.NamaPohon,
		Keterangan:  tematik.Keterangan,
		Indikators:  ConvertToIndikatorResponses(tematik.Indikator),
		SubTematiks: []pohonkinerja.SubtematikResponse{}, // Inisialisasi dengan array kosong
	}

	// Cek subtematik (level 1)
	if subTematiks := pohonMap[1][tematik.Id]; len(subTematiks) > 0 {
		var subTematikResponses []pohonkinerja.SubtematikResponse
		for _, subTematik := range subTematiks {
			subTematikResp := pohonkinerja.SubtematikResponse{
				Id:         subTematik.Id,
				Parent:     subTematik.Parent,
				Tema:       subTematik.NamaPohon,
				Keterangan: subTematik.Keterangan,
				Indikators: ConvertToIndikatorResponses(subTematik.Indikator),
				Strategics: []pohonkinerja.StrategicResponse{}, // Inisialisasi dengan array kosong
			}
			subTematikResponses = append(subTematikResponses, subTematikResp)
		}
		tematikResp.SubTematiks = subTematikResponses
	}

	return tematikResp
}
