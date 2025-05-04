import 'common.dart';
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

  CompetitorFull(this.id, this.fullName, this.birthDate, this.identity,
      this.bow, this.rank, this.region, this.federation, this.club);
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