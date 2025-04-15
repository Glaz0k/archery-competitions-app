import createEnum from "../../helper/createEnum";

const SportsRank = createEnum({
  MASTER_MERITED: { value: "merited_master", textValue: "ЗМС" },
  MASTER_INTERNATIONAL: { value: "master_international", textValue: "МСМК" },
  MASTER: { value: "master", textValue: "МС" },
  MASTER_CANDIDATE: { value: "candidate_for_master", textValue: "КМС" },
  CLASS_1: { value: "first_class", textValue: "1" },
  CLASS_2: { value: "second_class", textValue: "2" },
  CLASS_3: { value: "third_class", textValue: "3" },
});

export default SportsRank;
