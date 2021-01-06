package crops

//NASSCropMap is a map of Crops to NASS Crop ID #s
func NASSCropMap() map[string]Crop {
	m := make(map[string]Crop)

	m["1"] = BuildCrop(1, "Corn")
	m["2"] = BuildCrop(2,"Cotton")
	m["3"] = BuildCrop(3, "Rice")
	m["4"] = BuildCrop(4, "Sorghum")
	m["5"] = BuildCrop(5, "Soybeans")
	m["6"] = BuildCrop(6, "Sunflower")
	m["10"] = BuildCrop(10, "Peanuts")
	m["11"] = BuildCrop(11, "Tobacco")
	m["12"] = BuildCrop(12, "Sweet Corn")
	m["13"] = BuildCrop(13, "Pop or Orn Corn")
	m["14"] = BuildCrop(14, "Mint")
	m["21"] = BuildCrop(21, "Barley")
	m["22"] = BuildCrop(22, "Durum Wheat")
	m["23"] = BuildCrop(23, "Spring Wheat")
	m["24"] = BuildCrop(24, "Winter Wheat")
	m["25"] = BuildCrop(25, "Other Small Grains")
	m["26"] = BuildCrop(26,"Dbl Crop WinWht/Soybeans")
	m["27"] = BuildCrop(27, "Rye"}
	m["28"] = Crop{ID: 28, Name: "Oats"}
	m["29"] = Crop{ID: 29, Name: "Millet"}
	m["30"] = Crop{ID: 30, Name: "Speltz"}
	m["31"] = Crop{ID: 31, Name: "Canola"}
	m["32"] = Crop{ID: 32, Name: "Flaxseed"}
	m["33"] = Crop{ID: 33, Name: "Safflower"}
	m["34"] = Crop{ID: 34, Name: "Rape Seed"}
	m["35"] = Crop{ID: 35, Name: "Mustard"}
	m["36"] = Crop{ID: 36, Name: "Alfalfa"}
	m["37"] = Crop{ID: 37, Name: "Other Hay/Non Alfalfa"}
	m["38"] = Crop{ID: 38, Name: "Camelina"}
	m["39"] = Crop{ID: 39, Name: "Buckwheat"}
	m["41"] = Crop{ID: 41, Name: "Sugarbeets"}
	m["42"] = Crop{ID: 42, Name: "Dry Beans"}
	m["43"] = Crop{ID: 43, Name: "Potatoes"}
	m["44"] = Crop{ID: 44, Name: "Other Crops"}
	m["45"] = Crop{ID: 45, Name: "Sugarcane"}
	m["46"] = Crop{ID: 46, Name: "Sweet Potatoes"}
	m["47"] = Crop{ID: 47, Name: "Misc Vegs & Fruits"}
	m["48"] = Crop{ID: 48, Name: "Watermelons"}
	m["49"] = Crop{ID: 49, Name: "Onions"}
	m["50"] = Crop{ID: 50, Name: "Cucumbers"}
	m["51"] = Crop{ID: 51, Name: "Chick Peas"}
	m["52"] = Crop{ID: 52, Name: "Lentils"}
	m["53"] = Crop{ID: 53, Name: "Peas"}
	m["54"] = Crop{ID: 54, Name: "Tomatoes"}
	m["55"] = Crop{ID: 55, Name: "Caneberries"}
	m["56"] = Crop{ID: 56, Name: "Hops"}
	m["57"] = Crop{ID: 57, Name: "Herbs"}
	m["58"] = Crop{ID: 58, Name: "Clover/Wildflowers"}
	m["59"] = Crop{ID: 59, Name: "Sod/Grass Seed"}
	m["60"] = Crop{ID: 60, Name: "Switchgrass"}
	m["61"] = Crop{ID: 61, Name: "Fallow/Idle Cropland"}
	m["63"] = Crop{ID: 63, Name: "Forest"}
	m["64"] = Crop{ID: 64, Name: "Shrubland"}
	m["65"] = Crop{ID: 65, Name: "Barren"}
	m["66"] = Crop{ID: 66, Name: "Cherries"}
	m["67"] = Crop{ID: 67, Name: "Peaches"}
	m["68"] = Crop{ID: 68, Name: "Apples"}
	m["69"] = Crop{ID: 69, Name: "Grapes"}
	m["70"] = Crop{ID: 70, Name: "Christmas Trees"}
	m["71"] = Crop{ID: 71, Name: "Other Tree Crops"}
	m["72"] = Crop{ID: 72, Name: "Citrus"}
	m["74"] = Crop{ID: 74, Name: "Pecans"}
	m["75"] = Crop{ID: 75, Name: "Almonds"}
	m["76"] = Crop{ID: 76, Name: "Walnuts"}
	m["77"] = Crop{ID: 77, Name: "Pears"}
	m["92"] = Crop{ID: 92, Name: "Aquaculture"}
	m["152"] = Crop{ID: 152, Name: "Shrubland"}
	m["204"] = Crop{ID: 204, Name: "Pistachios"}
	m["205"] = Crop{ID: 205, Name: "Triticale"}
	m["206"] = Crop{ID: 206, Name: "Carrots"}
	m["207"] = Crop{ID: 207, Name: "Asparagus"}
	m["208"] = Crop{ID: 208, Name: "Garlic"}
	m["209"] = Crop{ID: 209, Name: "Cantaloupes"}
	m["210"] = Crop{ID: 210, Name: "Prunes"}
	m["211"] = Crop{ID: 211, Name: "Olives"}
	m["212"] = Crop{ID: 212, Name: "Oranges"}
	m["213"] = Crop{ID: 213, Name: "Honeydew Melons"}
	m["214"] = Crop{ID: 214, Name: "Broccoli"}
	m["215"] = Crop{ID: 215, Name: "Avocados"}
	m["216"] = Crop{ID: 216, Name: "Peppers"}
	m["217"] = Crop{ID: 217, Name: "Pomegranates"}
	m["218"] = Crop{ID: 218, Name: "Nectarines"}
	m["219"] = Crop{ID: 219, Name: "Greens"}
	m["220"] = Crop{ID: 220, Name: "Plums"}
	m["221"] = Crop{ID: 221, Name: "Strawberries"}
	m["222"] = Crop{ID: 222, Name: "Squash"}
	m["223"] = Crop{ID: 223, Name: "Apricots"}
	m["224"] = Crop{ID: 224, Name: "Vetch"}
	m["225"] = Crop{ID: 225, Name: "Dbl Crop WinWht/Corn"}
	m["226"] = Crop{ID: 226, Name: "Dbl Crop Oats/Corn"}
	m["227"] = Crop{ID: 227, Name: "Lettuce"}
	m["228"] = Crop{ID: 228, Name: "Dbl Crop Triticale/Corn"}
	m["229"] = Crop{ID: 229, Name: "Pumpkins"}
	m["230"] = Crop{ID: 230, Name: "Dbl Crop Lettuce/Durum Wheat"}
	m["231"] = Crop{ID: 231, Name: "Dbl Crop Lettuce/Cantaloupe"}
	m["232"] = Crop{ID: 232, Name: "Dbl Crop Lettuce/Cotton"}
	m["233"] = Crop{ID: 233, Name: "Dbl Crop Lettuce/Barley"}
	m["234"] = Crop{ID: 234, Name: "Dbl Crop Durum Wht/Sorghum"}
	m["235"] = Crop{ID: 235, Name: "Dbl Crop Barley/Sorghum"}
	m["236"] = Crop{ID: 236, Name: "Dbl Crop WinWht/Sorghum"}
	m["237"] = Crop{ID: 237, Name: "Dbl Crop Barley/Corn"}
	m["238"] = Crop{ID: 238, Name: "Dbl Crop WinWht/Cotton"}
	m["239"] = Crop{ID: 239, Name: "Dbl Crop Soybeans/Cotton"}
	m["240"] = Crop{ID: 240, Name: "Dbl Crop Soybeans/Oats"}
	m["241"] = Crop{ID: 241, Name: "Dbl Crop Corn/Soybeans"}
	m["242"] = Crop{ID: 242, Name: "Blueberries"}
	m["243"] = Crop{ID: 243, Name: "Cabbage"}
	m["244"] = Crop{ID: 244, Name: "Cauliflower"}
	m["245"] = Crop{ID: 245, Name: "Celery"}
	m["246"] = Crop{ID: 246, Name: "Radishes"}
	m["247"] = Crop{ID: 247, Name: "Turnips"}
	m["248"] = Crop{ID: 248, Name: "Eggplants"}
	m["249"] = Crop{ID: 249, Name: "Gourds"}
	m["250"] = Crop{ID: 250, Name: "Cranberries"}
	m["254"] = Crop{ID: 254, Name: "Dbl Crop Barley/Soybeans"}

	return m
}

/*
value,red,green,blue,category,opacity
1	255	211	0	Corn	255
2	255	38	38	Cotton	255
3	0	168	228	Rice	255
4	255	158	11	Sorghum	255
5	38	112	0	Soybeans	255
6	255	255	0	Sunflower	255
10	112	165	0	Peanuts	255
11	0	175	75	Tobacco	255
12	221	165	11	Sweet Corn	255
13	221	165	11	Pop or Orn Corn	255
14	126	211	255	Mint	255
21	226	0	124	Barley	255
22	137	98	84	Durum Wheat	255
23	216	181	107	Spring Wheat	255
24	165	112	0	Winter Wheat	255
25	214	158	188	Other Small Grains	255
26	112	112	0	Dbl Crop WinWht/Soybeans	255
27	172	0	124	Rye	255
28	160	89	137	Oats	255
29	112	0	73	Millet	255
30	214	158	188	Speltz	255
31	209	255	0	Canola	255
32	126	153	255	Flaxseed	255
33	214	214	0	Safflower	255
34	209	255	0	Rape Seed	255
35	0	175	75	Mustard	255
36	255	165	226	Alfalfa	255
37	165	242	140	Other Hay/Non Alfalfa	255
38	0	175	75	Camelina	255
39	214	158	188	Buckwheat	255
41	168	0	228	Sugarbeets	255
42	165	0	0	Dry Beans	255
43	112	38	0	Potatoes	255
44	0	175	75	Other Crops	255
45	177	126	255	Sugarcane	255
46	112	38	0	Sweet Potatoes	255
47	255	102	102	Misc Vegs & Fruits	255
48	255	102	102	Watermelons	255
49	255	204	102	Onions	255
50	255	102	102	Cucumbers	255
51	0	175	75	Chick Peas	255
52	0	221	175	Lentils	255
53	84	255	0	Peas	255
54	242	163	119	Tomatoes	255
55	255	102	102	Caneberries	255
56	0	175	75	Hops	255
57	126	211	255	Herbs	255
58	232	191	255	Clover/Wildflowers	255
59	175	255	221	Sod/Grass Seed	255
60	0	175	75	Switchgrass	255
61	191	191	119	Fallow/Idle Cropland	255
63	147	204	147	Forest	255
64	198	214	158	Shrubland	255
65	204	191	163	Barren	255
66	255	0	255	Cherries	255
67	255	142	170	Peaches	255
68	186	0	79	Apples	255
69	112	68	137	Grapes	255
70	0	119	119	Christmas Trees	255
71	177	154	112	Other Tree Crops	255
72	255	255	126	Citrus	255
74	181	112	91	Pecans	255
75	0	165	130	Almonds	255
76	233	214	175	Walnuts	255
77	177	154	112	Pears	255
92	0	255	255	Aquaculture	255
152	198	214	158	Shrubland	255
204	0	255	140	Pistachios	255
205	214	158	188	Triticale	255
206	255	102	102	Carrots	255
207	255	102	102	Asparagus	255
208	255	102	102	Garlic	255
209	255	102	102	Cantaloupes	255
210	255	142	170	Prunes	255
211	51	73	51	Olives	255
212	228	112	38	Oranges	255
213	255	102	102	Honeydew Melons	255
214	255	102	102	Broccoli	255
215	102	153	76	Avocados	255
216	255	102	102	Peppers	255
217	177	154	112	Pomegranates	255
218	255	142	170	Nectarines	255
219	255	102	102	Greens	255
220	255	142	170	Plums	255
221	255	102	102	Strawberries	255
222	255	102	102	Squash	255
223	255	142	170	Apricots	255
224	0	175	75	Vetch	255
225	255	211	0	Dbl Crop WinWht/Corn	255
226	255	211	0	Dbl Crop Oats/Corn	255
227	255	102	102	Lettuce	255
228	255	210	0	Dbl Crop Triticale/Corn	255
229	255	102	102	Pumpkins	255
230	137	98	84	Dbl Crop Lettuce/Durum Wht	255
231	255	102	102	Dbl Crop Lettuce/Cantaloupe	255
232	255	38	38	Dbl Crop Lettuce/Cotton	255
233	226	0	124	Dbl Crop Lettuce/Barley	255
234	255	158	11	Dbl Crop Durum Wht/Sorghum	255
235	255	158	11	Dbl Crop Barley/Sorghum	255
236	165	112	0	Dbl Crop WinWht/Sorghum	255
237	255	211	0	Dbl Crop Barley/Corn	255
238	165	112	0	Dbl Crop WinWht/Cotton	255
239	38	112	0	Dbl Crop Soybeans/Cotton	255
240	38	112	0	Dbl Crop Soybeans/Oats	255
241	255	211	0	Dbl Crop Corn/Soybeans	255
242	0	0	153	Blueberries	255
243	255	102	102	Cabbage	255
244	255	102	102	Cauliflower	255
245	255	102	102	Celery	255
246	255	102	102	Radishes	255
247	255	102	102	Turnips	255
248	255	102	102	Eggplants	255
249	255	102	102	Gourds	255
250	255	102	102	Cranberries	255
254	38	112	0	Dbl Crop Barley/Soybeans	255
*/
