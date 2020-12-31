package crops
func NASSCropMap() map[string]Crop {
	m := make(map[string]Crop)
	m["1"] = Crop{ID:1, Name:"Corn"}
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
