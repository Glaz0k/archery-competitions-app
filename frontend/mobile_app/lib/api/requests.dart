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

// "login": <string>,
// "password": <string>
class Credentials {
  final String login;
  final String password;

  Credentials(this.login, this.password);
}

// "range_ordinal": <number>,
// "shots": <[ <shot> ] | null>
class ChangeRange {
  final int rangeOrdinal;
  final List<Shot>? shots;

  ChangeRange(this.rangeOrdinal, this.shots);
}

// "score": <string>,
// "priority": <bool | null>
class ChangeShootOut {
  final String score;
  final bool? priority;

  ChangeShootOut(this.score, this.priority);
}
