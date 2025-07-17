import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lunary_app/app/lunary_app.dart';

void main() {
  testWidgets('LUNARY app smoke test', (WidgetTester tester) async {
    // Build our app and trigger a frame.
    await tester.pumpWidget(
      ProviderScope(
        child: LunaryApp(),
      ),
    );

    // Verify that app loads without crashing
    expect(find.text('LUNARY'), findsOneWidget);
  });
}
