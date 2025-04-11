const mockTime = 2000;

export async function getCups() {
  await new Promise((res) => setTimeout(res, mockTime));
  const data = [
    {
      id: 0,
      title: "Gink",
      address: "c. Linwood, st. Douglass Street 27",
      season: null,
    },
    {
      id: 1,
      title: "Blanet",
      address: null,
      season: "2023/2024",
    },
    {
      id: 2,
      title: "Ecraze",
      address: "c. Century, st. Roosevelt Place 28",
      season: null,
    },
    {
      id: 3,
      title: "Gracker",
      address: null,
      season: "2023/2024",
    },
    {
      id: 4,
      title: "Zytrex",
      address: null,
      season: "2023/2024",
    },
    {
      id: 5,
      title: "Nebulean",
      address: "c. Winfred, st. Glen Street 39",
      season: "2023/2024",
    },
    {
      id: 6,
      title: "Zialactic",
      address: "c. Thornport, st. Rockaway Avenue 62",
      season: null,
    },
    {
      id: 7,
      title: "Ozean",
      address: null,
      season: "2023/2024",
    },
    {
      id: 8,
      title: "Ecolight",
      address: null,
      season: null,
    },
    {
      id: 9,
      title: "Quantalia",
      address: "c. Abrams, st. Downing Street 66",
      season: null,
    },
    {
      id: 10,
      title: "Zaggle",
      address: null,
      season: null,
    },
    {
      id: 11,
      title: "Earbang",
      address: null,
      season: "2023/2024",
    },
    {
      id: 12,
      title: "Inventure",
      address: null,
      season: "2023/2024",
    },
    {
      id: 13,
      title: "Opticom",
      address: "c. Mathews, st. Allen Avenue 67",
      season: null,
    },
    {
      id: 14,
      title: "Zytrax",
      address: null,
      season: null,
    },
    {
      id: 15,
      title: "Isoswitch",
      address: "c. Alleghenyville, st. Hendrix Street 28",
      season: null,
    },
    {
      id: 16,
      title: "Quilch",
      address: null,
      season: null,
    },
    {
      id: 17,
      title: "Chillium",
      address: null,
      season: "2023/2024",
    },
    {
      id: 18,
      title: "Mediot",
      address: "c. Accoville, st. Reeve Place 70",
      season: null,
    },
    {
      id: 19,
      title: "Equitox",
      address: null,
      season: null,
    },
  ];
  return data.slice();
}

export async function postCup(addCupRequestBody) {
  await new Promise((res) => setTimeout(res, mockTime));
  if (addCupRequestBody.title == "") {
    throw {
      code: 400,
      error: "INVALID PARAMETERS",
    };
  }
  if (addCupRequestBody.address == "") {
    addCupRequestBody.address = null;
  }
  if (addCupRequestBody.season == "") {
    addCupRequestBody.season = null;
  }
  const cup = {
    id: Math.floor(Math.random() * 10000),
    title: addCupRequestBody.title,
    address: addCupRequestBody.address,
    season: addCupRequestBody.season,
  };
  return cup;
}

export async function deleteCup(_cupId) {
  await new Promise((res) => setTimeout(res, mockTime));
  return true;
}
