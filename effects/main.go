// The EffectFactory allows to do magic effects on your actors.
// Effects like scale in/out, shake, fade in/out are supported.
package effect

// public enum EffectType {
//   None,
//   ScaleIn, ScaleOut, ScaleInOut, ScaleTo, ScaleToBack,
//   ShakeIn, ShakeOut, ShakeInOut, ShakeTo, ShakeToBack,
//   FadeIn, FadeOut, FadeInOut, FadeTo, FadeToBack,
//   ScaleToThenFadeOut, ScaleAndFadeOut, ScaleAndFadeIn,

//   SlideRight, SlideLeft, SlideUp, SlideDown,
//   PatrolX, PatrolY;
// }

// enum EffectDuration {
//   Once,
//   OnceToAndBack,
//   Looping,
//   LoopingToAndBack
// }

// public class Effect {

//   public static void createEffect(Actor actor, EffectType effectType, float value,
//       float duration, InterpolationType type){
//     if(actor == null)
//       return;
//     Interpolation interp = InterpolationType.getInterpolation(type);
//     float x = 0;//actor.getX();
//     float y = 0;//actor.getY();
//     switch(effectType){
//       case ScaleToThenFadeOut:
//         actor.addAction(Actions.sequence(Actions.scaleTo(value, value, duration, interp)));
//         break;
//       case PatrolX:
//         actor.addAction(Actions.forever(Actions.sequence(Actions.moveBy(value, 0, duration, interp),
//             Actions.moveBy(-value, 0, duration, interp))));
//         break;

//       case PatrolY:
//         actor.addAction(Actions.forever(Actions.sequence(Actions.moveBy(0, value, duration, interp),
//             Actions.moveBy(0, -value, duration, interp))));
//         break;

//       case SlideLeft:
//         actor.setPosition(999, y);
//         actor.addAction(Actions.moveTo(x, y, duration, interp));
//         break;
//       case SlideRight:
//         actor.setPosition(-999, y);
//         actor.addAction(Actions.moveTo(x, y, duration, interp));
//         break;
//       case SlideUp:
//         actor.setPosition(x, -999);
//         actor.addAction(Actions.moveTo(x, y, duration, interp));
//         break;
//       case SlideDown:
//         actor.setPosition(x, 999);
//         actor.addAction(Actions.moveTo(x, y, duration, interp));
//         break;
//       case FadeIn:
//         Color color = actor.getColor();
//         color.a = 0f;
//         actor.setColor(color);
//         actor.addAction(Actions.fadeIn(duration, interp));
//         break;
//       case FadeOut:
//         Color color2 = actor.getColor();
//         color2.a = 1f;
//         actor.setColor(color2);
//         actor.addAction(Actions.fadeOut(duration, interp));
//         break;
//       case FadeInOut:
//         actor.addAction(fadeInOut(value, duration, interp));
//         break;
//       case ScaleIn:
//         actor.setScale(0, 0);
//         actor.addAction(Actions.scaleTo(1, 1, duration, interp));
//         break;
//       case ScaleOut:
//         actor.setScale(1, 1);
//         actor.addAction(Actions.scaleTo(0, 0, duration, interp));
//         break;
//       case ScaleInOut:
//         actor.addAction(scaleInOut(value, duration, interp));
//         break;
//       case None:
//         break;
//       default:
//         break;
//     }
//   }

//   public static void createEffect(ImageJson imageJson){
//     createEffect(imageJson, imageJson.effectType, imageJson.effectValue, imageJson.effectDuration,
//         imageJson.interpolationType);
//   }

//   public static Action scaleInOut(float value, float duration, Interpolation interp){
//     return Actions.sequence(Actions.scaleTo(value, value, duration, interp), Actions.scaleTo(1, 1, duration, interp));
//   }

//   public static Action shakeInOut(float value, float duration, Interpolation interp){
//     return Actions.sequence(Actions.rotateTo(value, duration, interp), Actions.rotateTo(-value, duration, interp),
//         Actions.rotateTo(0, duration, interp));
//   }

//   public static Action fadeInOut(float value, float duration, Interpolation interp){
//     return Actions.sequence(Actions.fadeIn(duration, interp), Actions.fadeOut(duration, interp));
//   }

//   public static Action fadeOutIn(float value, float duration, Interpolation interp){
//     return Actions.sequence(Actions.fadeOut(duration, interp), Actions.fadeIn(duration, interp));
//   }

//   /**
//    * Scale effect, Shake effect (SC, SHK)
//    * */
//   public static void create_SC_SHK(Actor actor, float scaleRatioX,
//       float scaleRatioY, float shakeAngle, float originalAngle,
//       float duration, final boolean isRemoveActor) {
//     if (actor != null) {
//       actor.addAction(Actions.sequence(
//           Actions.scaleTo(scaleRatioX, scaleRatioY, duration),
//           Actions.rotateTo(shakeAngle, duration),
//           Actions.rotateTo(-shakeAngle, duration),
//           Actions.rotateTo(originalAngle, duration)));
//     }
//   }

//   //
//   // TRIPLE EFFECTS (Sequence - Chain reaction)
//   // ################################################################
//   /**
//    * Scale effect, Back To Normal, Fade Out (SC, BTN, FO)
//    * */
//   public static void create_SC_BTN_FO(Actor actor, float scaleRatioX,
//       float scaleRatioY, float duration, float delayBeforeFadeOut,
//       final boolean isRemoveActor) {
//     if (actor != null) {
//       actor.addAction(Actions.sequence(
//           Actions.scaleTo(scaleRatioX, scaleRatioY, duration),
//           Actions.scaleTo(1, 1, duration),
//           Actions.delay(delayBeforeFadeOut),
//           Actions.fadeOut(duration)));
//     }
//   }

//   /**
//    * Scale effect, Shake effect, Back To Normal (SC, SHK, BTN)
//    * */
//   public static void create_SC_SHK_BTN(Actor actor, float scaleRatioX,
//       float scaleRatioY, float shakeAngle, float originalAngle,
//       float duration, final boolean isRemoveActor) {
//     if (actor != null) {
//       actor.addAction(Actions.sequence(
//           Actions.scaleTo(scaleRatioX, scaleRatioY, duration),
//           Actions.rotateTo(shakeAngle, duration),
//           Actions.rotateTo(-shakeAngle, duration),
//           Actions.rotateTo(originalAngle, duration),
//           Actions.scaleTo(1, 1, duration)));
//     }
//   }

//   static Interpolation[] interpolationsValue = {
//     Interpolation.bounce, Interpolation.bounceIn,  Interpolation.bounceOut,
//     Interpolation.circle,  Interpolation.circleIn,  Interpolation.circleOut,
//     Interpolation.elastic,  Interpolation.elasticIn,  Interpolation.elasticOut,
//     Interpolation.exp10,  Interpolation.exp10In,  Interpolation.exp10Out,
//     Interpolation.exp5,  Interpolation.exp5In,  Interpolation.exp5Out,
//     Interpolation.linear,  Interpolation.fade,
//     Interpolation.pow2,  Interpolation.pow2In,  Interpolation.pow2Out,
//     Interpolation.pow3,  Interpolation.pow3In,  Interpolation.pow3Out,
//     Interpolation.pow4,  Interpolation.pow4In,  Interpolation.pow4Out,
//     Interpolation.pow5,  Interpolation.pow5In,  Interpolation.pow5Out,
//     Interpolation.sine,  Interpolation.sineIn,  Interpolation.sineOut,
//     Interpolation.swing,  Interpolation.swingIn,  Interpolation.swingOut,
//   };
// }
