import createEnum from "../helper/createEnum";

const GroupGender = createEnum({
  MALE: { value: "male", textValue: "Мужчины" },
  FEMALE: { value: "female", textValue: "Женщины" },
  UNITED: { value: "united", textValue: "Объединённый" },
});

export default GroupGender;
