import 'dart:developer';

import 'package:flutter/material.dart';
import 'package:mobile_app/model/user_model.dart';
import 'package:mobile_app/page/range_input_page.dart';
import 'package:mobile_app/page/widgets/qualification_table.dart';
import 'package:provider/provider.dart';

import '../api/api.dart';
import '../api/responses.dart';
import '../model/range_model.dart';

class QualificationPage extends StatefulWidget {
  final int groupId;

  const QualificationPage({super.key, required this.groupId});

  @override
  State<QualificationPage> createState() => _QualificationPageState();
}

class _QualificationPageState extends State<QualificationPage> {
  Section? _section;
  List<QualificationRoundFull>? _rounds;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    var api = context.read<Api>();
    int userId = context.read<UserModel>().getId()!;
    QualificationTable table = await api.getIndividualGroupQualificationTable(
      widget.groupId,
    );
    try {
      var section = table.sections.firstWhere(
        (section) => section.competitor.id == userId,
      );
      setState(() {
        _section = section;
      });
      var futures =
          section.rounds
              .map(
                (round) => api.getQualificationSectionsRound(
                  section.id,
                  round.roundOrdinal,
                ),
              )
              .toList();
      log("Загружаем раунды $futures");
      var rounds = await Future.wait(futures);
      log("Загрузили $rounds");
      rounds.sort((a, b) => a.roundOrdinal.compareTo(b.roundOrdinal));
      setState(() {
        _rounds = rounds;
      });
    } on StateError {
      setState(() {
        _section = null;
        _rounds = null;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    var rounds = _rounds;
    return RefreshIndicator(
      onRefresh: _loadData,
      child: SingleChildScrollView(
        physics: AlwaysScrollableScrollPhysics(),
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              MyQualificationTable(section: _section),
              if (rounds != null)
                for (var round in rounds) RoundCard(round: round),
            ],
          ),
        ),
      ),
    );
  }
}

class RoundCard extends StatelessWidget {
  final QualificationRoundFull round;

  const RoundCard({super.key, required this.round});

  @override
  Widget build(BuildContext context) {
    var api = context.watch<Api>();
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: [
            Column(
              children: [
                Text(
                  "Раунд ${round.roundOrdinal}",
                  style: Theme.of(context).textTheme.headlineSmall,
                ),
                if (round.rangeGroup.totalScore != null)
                  Text("Счёт: ${round.rangeGroup.totalScore}"),
              ],
            ),
            if (round.isActive)
              FilledButton(
                onPressed:
                    () => Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder:
                            (context) => ChangeNotifierProvider(
                              child: SeriesPage(),
                              create:
                                  (context) => RangeModel(
                                    round.rangeGroup,
                                    getRangeGroup:
                                        () => api
                                            .getQualificationSectionsRoundsRanges(
                                              round.sectionId,
                                              round.roundOrdinal,
                                            ),
                                    putRange:
                                        (request) => api
                                            .putQualificationSectionsRoundsRange(
                                              round.sectionId,
                                              round.roundOrdinal,
                                              request,
                                            ),
                                    endRange:
                                        (idx) => api
                                            .endQualificationSectionsRoundsRange(
                                              round.sectionId,
                                              round.roundOrdinal,
                                              idx,
                                            ),
                                  ),
                            ),
                      ),
                    ),
                child: Text("Перейти"),
              ),
          ],
        ),
      ),
    );
  }
}
