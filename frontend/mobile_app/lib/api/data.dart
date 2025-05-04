// "id": <number>,
// "full_name": <string>,
// "birth_date": <date ISO 8601>,
// "identity": <gender>,
// "bow": <bow_class | null>,
// "rank": <sports_rank | null>,
// "region": <string | null>,
// "federation": <string | null>,
// "club": <string | null>
class CompetitorFull {
  final int id;
  final String fullName;
  final DateTime birthDate;
  final Gender identity;
  final BowClass? bow;
  final SportsRank? rank;
  final String? region;
  final String? federation;
  final String? club;

  CompetitorFull(
    this.id,
      this.fullName,
      this.birthDate,
      this.identity,
      this.bow,
      this.rank,
      this.region,
      this.federation,
      this.club
      );
}


// "competition_id": <number>,
// "competitor": <competitor_full>,
// "is_active": <bool>,
// "created_at": <YYYY-MM-DDThh:mm:ss±hh ISO 8601>
class CompetitorCompetitionDetail {
  final int competitionId;
  final CompetitorFull competitor;
  final bool isActive;
  final DateTime createdAt;

  CompetitorCompetitionDetail(this.competitionId, this.competitor,
      this.isActive, this.createdAt);
}

// "id": <number>,
// "competition_id": <number>,
// "bow": <bow_class>,
// "identity": <gender | null>,
// "state": <group_state>
class IndividualGroup {
  final int id;
  final int competitionId;
  final BowClass bow;
  final Gender? identity;
  final GroupState state;

  IndividualGroup(this.id, this.competitionId, this.bow, this.identity,
      this.state);

}

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