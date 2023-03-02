package structureprovider

/*
	func TestGPKByFips(t *testing.T) {
		filepath := "/workspaces/Go_Consequences/data/nsiv2_11.gpkg"
		nsp, _ := InitGPK(filepath, "nsi")
		fmt.Println(nsp.FilePath)
		d := hazards.DepthEvent{}
		d.SetDepth(2.4)
		nsp.ByFips("11001001600", func(s consequences.Receptor) {
			r, _ := s.Compute(d)
			n, err := r.Fetch("fd_id")
			if err != nil {
				panic(err)
			}
			dam, err := r.Fetch("structure damage")
			if err != nil {
				panic(err)
			}
			name, nok := n.(string)
			if !nok {
				panic("n was not ok.")
			}
			damage, damok := dam.(float64)
			if !damok {
				panic("dam was not ok.")
			}
			fmt.Println(name + " had " + fmt.Sprint(damage))
		})
	}
*/
/*func TestGPKByBBox(t *testing.T) {
	filepath := "/workspaces/Go_Consequences/data/nsiv2_11.gpkg"
	nsp, _ := InitGPK(filepath, "nsi")
	fmt.Println(nsp.FilePath)
	d := hazards.DepthEvent{}
	d.SetDepth(2.4)
	bbox := make([]float64, 4)
	bbox[0] = -79.00 //upper left x
	bbox[1] = 39.00  //upper left y
	bbox[2] = -76.00 //lower right x
	bbox[3] = 38.00  //lower right y
	gbbx := geography.BBox{Bbox: bbox}
	nsp.ByBbox(gbbx, func(s consequences.Receptor) {
		r, _ := s.Compute(d)
		b, _ := json.Marshal(r)
		fmt.Println(string(b))
	})
}*/
