// merited_master
// master_international
// master
// candidate_for_master
// first_class
// second_class
// third_class
enum SportsRank {
  meritedMaster,
  masterInternational,
  master,
  candidateForMaster,
  firstClass,
  secondClass,
  thirdClass,
}

// male
// female
enum Gender { male, female }

// classic
// block
// classic_newbie
// 3D_classic
// 3D_compound
// 3D_long
// peripheral
// peripheral_with_ring
enum BowClass {
  classic,
  block,
  classicNewbie,
  classic3D,
  compound3D,
  long3D,
  peripheral,
  peripheralWithRing,
}

// created
// qualification_start
// qualification_end
// quarterfinal_start
// semifinal_start
// final_start
// completed
enum GroupState {
  created,
  qualificationStart,
  qualificationEnd,
  quarterfinalStart,
  semifinalStart,
  finalStart,
  completed,
}

// I
// II
// III
// F
enum CompetitionStage { I, II, III, F }

// ongoing
// top_win
// bot_win
enum SparringState { ongoing, topWin, botWin }

// "1-10"
// "6-10"
enum RangeType { one2ten, six2ten }

// "shot_ordinal": <number>,
// "score": <string | null>
class Shot {
  int shotOrdinal;
  String? score;

  Shot(this.shotOrdinal, this.score);
}
