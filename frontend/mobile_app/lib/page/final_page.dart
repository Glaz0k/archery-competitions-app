import 'package:flutter/material.dart';
import 'package:mobile_app/api/responses.dart';
import 'package:provider/provider.dart';

import '../api/api.dart';

class FinalPage extends StatefulWidget {
  final int groupId;

  const FinalPage({super.key, required this.groupId});

  @override
  State<StatefulWidget> createState() => _FinalPageState();
}

class _FinalPageState extends State<FinalPage> {
  FinalGrid? _grid;

  @override
  void initState() {
    super.initState();
    var api = context.read<Api>();
    api
        .getIndividualGroupFinalGrid(widget.groupId)
        .then(
          (grid) => setState(() {
            _grid = grid;
          }),
        );
  }

  @override
  Widget build(BuildContext context) {
    var api = context.watch<Api>();
    var grid = _grid;
    return RefreshIndicator(
      child: ListView(scrollDirection: Axis.vertical, children: [
        if (grid?.fina1?.sparringGold != null)
          FinalCard(title: "Финал", sparringList: [grid?.fina1?.sparringGold]),
        if (grid?.fina1?.sparringBronze != null)
          FinalCard(title: "Финал за 3-е место", sparringList: [grid?.fina1?.sparringBronze]),
        if (grid?.semifinal != null) 
          FinalCard(title: "Полуфинал", sparringList: [
            grid?.semifinal?.sparring5,
            grid?.semifinal?.sparring6,
          ]),
        if (grid?.quarterfinal != null) 
          FinalCard(title: "Четвертьфинал", sparringList: [
            grid?.quarterfinal.sparring1,
            grid?.quarterfinal.sparring2,
            grid?.quarterfinal.sparring3,
            grid?.quarterfinal.sparring4,
          ])
      ]),
      onRefresh: () async {
        grid = await api.getIndividualGroupFinalGrid(widget.groupId);
        setState(() {
          _grid = grid;
        });
      },
    );
  }
}

class FinalCard extends StatelessWidget {
  final List<Sparring?> sparringList;
  final String title;

  const FinalCard({super.key, required this.title, required this.sparringList});

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Column(
        children: [
          Text(title),
          Divider(),
          for (var sparring in sparringList)
            if (sparring != null) SparringCard(sparring: sparring),
        ],
      ),
    );
  }
}

class SparringCard extends StatelessWidget {
  final Sparring sparring;

  const SparringCard({super.key, required this.sparring});

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(sparring.topPlace?.competitor.fullName ?? "Нет соперника"),
        Text("vs"),
        Text(sparring.botPlace?.competitor.fullName ?? "Нет соперника"),
      ],
    );
  }
}
