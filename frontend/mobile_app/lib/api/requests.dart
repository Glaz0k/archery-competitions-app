import 'common.dart';

// "full_name": <string>,
// "birth_date": <date ISO 8601>,
// "identity": <gender>,
// "bow": <bow_class | null>,
// "rank": <sports_rank | null>,
// "region": <string | null>,
// "federation": <string | null>,
// "club": <string | null>
class ChangeCompetitor {
  final String fullName;
  final DateTime birthDate;
  final Gender identity;
  final BowClass? bow;
  final SportsRank? rank;
  final String? region;
  final String? federation;
  final String? club;

  ChangeCompetitor(
    this.fullName,
    this.birthDate,
    this.identity,
    this.bow,
    this.rank,
    this.region,
    this.federation,
    this.club,
  );
}
